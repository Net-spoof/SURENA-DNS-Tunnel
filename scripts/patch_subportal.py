#!/usr/bin/env python3
"""Fail-closed DNSFabric integration for pinned WhiteDNS panel-v0.8.2.

Patches only the panel HTTP/UI layer; does not change the data model,
credentials, tunnel engine or existing /sub/ subscription endpoint.
"""
from pathlib import Path
import shutil
import sys

if len(sys.argv) != 2:
    raise SystemExit("Usage: python3 scripts/patch_subportal.py /path/to/whitedns-source")

root = Path(sys.argv[1]).resolve()
assets = Path(__file__).resolve().parent.parent / "patches"
http = root / "internal/panel/http.go"
subs = root / "internal/panel/subscriptions.go"
app = root / "internal/panel/web/app.js"
html = root / "internal/panel/web/index.html"
go_portal = root / "internal/panel/dnsfabric_subportal.go"

def replace_exact(text, before, after, count=1):
    actual = text.count(before)
    if actual != count:
        raise RuntimeError(f"Expected {count} copies of {before[:70]!r}, got {actual}. Unsupported upstream.")
    return text.replace(before, after)

if go_portal.exists():
    raise SystemExit("Refusing to apply DNSFabric portal twice")

# Stage all changes in memory; a failed validation changes no upstream file.
h = http.read_text()
old = 'if strings.HasPrefix(r.URL.Path, "/sub/") {'
if h.count(old) != 2:
    raise SystemExit("Upstream HTTP route structure differs from pinned version")
h = h.replace(old,
    'if strings.HasPrefix(r.URL.Path, "/sub/") || strings.HasPrefix(r.URL.Path, "/s/") || strings.HasPrefix(r.URL.Path, "/u/") {',
    1)
route = '''if strings.HasPrefix(r.URL.Path, "/sub/") {
\t\th.serveSubscription(w, r)
\t\treturn
\t}'''
route_new = route + '''
\tif strings.HasPrefix(r.URL.Path, "/s/") {
\t\t// The short URL reuses the official validation, headers and body.
\t\tcopy := r.Clone(r.Context())
\t\tlocation := *r.URL
\t\tlocation.Path = "/sub/" + strings.TrimPrefix(r.URL.Path, "/s/")
\t\tcopy.URL = &location
\t\th.serveSubscription(w, copy)
\t\treturn
\t}
\tif strings.HasPrefix(r.URL.Path, "/u/") {
\t\th.serveDNSFabricPortal(w, r)
\t\treturn
\t}'''
h = replace_exact(h, route, route_new)

s = subs.read_text()
s = replace_exact(s,
    '\tlink := h.subscriptionURL(r, credential)\n',
    '\tlink := h.subscriptionURL(r, credential)\n'
    '\tshort := strings.Replace(link, "/sub/", "/s/", 1)\n'
    '\tportal := strings.Replace(link, "/sub/", "/u/", 1)\n')
s = replace_exact(s, "qrcode.Encode(link, qrcode.Medium, 512)",
    "qrcode.Encode(short, qrcode.Medium, 512)")
s = replace_exact(s, '\turls := []string{link}\n',
    '\turls := []string{link}\n\tshortURLs := []string{short}\n')
s = replace_exact(s, '\t\t\t\turls = append(urls, alternative)\n',
    '\t\t\t\turls = append(urls, alternative)\n'
    '\t\t\t\tshortURLs = append(shortURLs, strings.Replace(alternative, "/sub/", "/s/", 1))\n')
s = replace_exact(s,
    'jsonReply(w, 200, map[string]any{"url": link, "urls": urls, "update_interval_hours": 12})',
    'jsonReply(w, 200, map[string]any{"url": link, "urls": urls, '
    '"short_url": short, "short_urls": shortURLs, "portal_url": portal, '
    '"update_interval_hours": 12})')

js = app.read_text()
js = replace_exact(js, "subscriptionLink='',subscriptionURLs=[];",
    "subscriptionLink='',subscriptionURLs=[],subscriptionPortalURL='';")
js = replace_exact(js,
    "subscriptionLink=sub.url;subscriptionURLs=sub.urls||[sub.url];",
    "subscriptionLink=sub.short_url||sub.url;"
    "subscriptionURLs=sub.short_urls||sub.urls||[sub.url];"
    "subscriptionPortalURL=sub.portal_url||'';")
js = replace_exact(js,
    "subscriptionLink=v.url;subscriptionURLs=v.urls||[v.url];",
    "subscriptionLink=v.short_url||v.url;"
    "subscriptionURLs=v.short_urls||v.urls||[v.url];"
    "subscriptionPortalURL=v.portal_url||'';")
js = replace_exact(js,
    "$('rotate-subscription').hidden=!sub;",
    "$('share-portal').hidden=!sub||!subscriptionPortalURL;"
    "$('rotate-subscription').hidden=!sub;")
js = replace_exact(js,
    "$('share-sub-mode').onclick=()=>showShareMode('subscription');",
    "$('share-sub-mode').onclick=()=>showShareMode('subscription');\n"
    "$('share-portal').onclick=()=>{if(subscriptionPortalURL)"
    "window.open(subscriptionPortalURL,'_blank','noopener,noreferrer');};")

page = html.read_text()
page = replace_exact(page,
    '<button id="rotate-subscription"',
    '<button id="share-portal" class="quiet" type="button" hidden>'
    'نمایش صفحه اشتراک · Subscription portal</button>'
    '<button id="rotate-subscription"')

testfile = assets / "dnsfabric_subportal_test.go"
source = assets / "dnsfabric_subportal.go"
for p in (source, testfile):
    if not p.is_file():
        raise SystemExit(f"Missing patch payload: {p}")

for p, body in ((http, h), (subs, s), (app, js), (html, page)):
    p.write_text(body)
shutil.copyfile(source, go_portal)
shutil.copyfile(testfile, root / "internal/panel/dnsfabric_subportal_test.go")
print("PASS: DNSFabric portal integrated into WhiteDNS pinned source, /sub/ preserved")

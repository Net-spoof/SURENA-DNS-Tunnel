package panel

import (
  "fmt"
  "html"
  "net/http"
  "strings"
  "time"
  qrcode "github.com/skip2/go-qrcode"
)

// DNSFabric subscriber view; token remains the official WhiteDNS credential.
func (h *HTTP) serveDNSFabricPortal(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" && r.Method != "HEAD" {
    w.Header().Set("Allow","GET, HEAD")
    w.WriteHeader(405)
    return
  }
  token := strings.TrimPrefix(r.URL.Path, "/u/")
  c, _, err := h.store.Subscription(token)
  if err != nil { http.NotFound(w,r);return }
  status:=c.Status(time.Now())
  base:=h.subscriptionURL(r, token)
  url:=strings.Replace(base, "/sub/", "/s/",1)
  w.Header().Set("Cache-Control","no-store, private")
  w.Header().Set("Referrer-Policy","no-referrer")
  w.Header().Set("X-Robots-Tag","noindex, nofollow")
  w.Header().Set("X-Content-Type-Options","nosniff")
  w.Header().Set("X-Frame-Options","DENY")
  w.Header().Set("Content-Security-Policy","default-src 'none'; img-src 'self'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; frame-ancestors 'none'; base-uri 'none'")
  if r.URL.Query().Get("qr")=="1" {
    if status!="active" { http.NotFound(w,r);return }
    png,e:=qrcode.Encode(url,qrcode.Medium,256)
    if e!=nil {http.Error(w,"QR error",500);return}
    w.Header().Set("Content-Type","image/png")
    if r.Method=="GET" {_,_=w.Write(png)}
    return
  }
  name:=html.EscapeString(c.Name)
  used:=float64(c.Used())/1073741824
  quotaText,remainText,ratio:="Unlimited","Unlimited",float64(0)
  if c.QuotaBytes>0 {
    quotaText=fmt.Sprintf("%.2f GiB",float64(c.QuotaBytes)/1073741824)
    remainText=fmt.Sprintf("%.2f GiB",float64(max(int64(0),c.QuotaBytes-c.Used()))/1073741824)
    ratio=100*float64(c.Used())/float64(c.QuotaBytes)
    if ratio>100 {ratio=100}
  }
  expiry:="Unlimited"
  if c.ExpiresAt>0 {expiry=time.UnixMilli(c.ExpiresAt).UTC().Format("2006-01-02")}
  action:="<p class='inactive'>این اشتراک غیرفعال یا منقضی شده است. برای تمدید با ارائه‌دهنده تماس بگیرید. · Subscription inactive.</p>"
  if status=="active" {
    action=fmt.Sprintf("<div class='actions'><a class='btn' href='%s' rel='noreferrer'>دریافت اشتراک · Import</a><button id='copy' type='button'>کپی لینک · Copy</button></div><input id='subscription' readonly dir='ltr' value='%s' aria-label='Subscription URL'><div class='qr'><img width='160' height='160' alt='QR Code' src='/u/%s?qr=1'><p>اسکن QR یا افزودن لینک به WhiteDNS</p></div>",html.EscapeString(url),html.EscapeString(url),html.EscapeString(token))
  }
  body:=fmt.Sprintf(`<!doctype html><html lang="fa" dir="rtl"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>DNSFabric | پنل اشتراک</title>
<style>*{box-sizing:border-box}body{margin:0;background:radial-gradient(at top left,#163a58,#0c121e 65%%);color:#ebf7fa;font:16px system-ui;min-height:100vh;display:grid;place-items:center;padding:16px}.card{max-width:660px;width:100%%;background:#172535;border:1px solid #365065;border-radius:22px;padding:28px;box-shadow:0 25px 60px #0008}header{color:#70e3d1;font-size:26px;font-weight:800}small{color:#a7b9c9}h1{overflow-wrap:anywhere}section{display:grid;grid-template-columns:repeat(2,1fr);gap:12px}.metric{background:#243649;border-radius:14px;padding:16px}b{display:block;margin-top:8px}progress{width:100%%;accent-color:#70e3d1}.actions{display:flex;gap:10px;flex-wrap:wrap;margin:22px 0 14px}.btn,button{display:inline-block;background:#60e3ca;color:#061a1e;padding:14px 18px;border-radius:12px;text-decoration:none;font-weight:700;border:0;font:inherit;cursor:pointer}button{background:#28495d;color:#e3f8fa}.state{color:#b3c6d6}.state b{color:#77e5c8}.qr{text-align:center;margin:22px 0 4px;padding-top:18px;border-top:1px solid #ffffff20}.qr img{background:white;padding:6px;border-radius:12px}.qr p{color:#b3c5d2;font-size:13px}input{width:100%%;direction:ltr;background:#0d1826;color:#cde2ed;border:1px solid #4a6477;padding:13px;border-radius:12px;font:12px monospace}.inactive{color:#f1bc96;line-height:1.8}footer{margin-top:26px;font-size:12px;color:#8fa7b8}@media(max-width:450px){section{grid-template-columns:1fr}}</style></head><body><article class="card"><header>DNSFabric</header><small>WHITE DNS · SUBSCRIPTION PORTAL</small><h1>%s</h1><p class="state">وضعیت سرویس · Status: <b>%s</b></p><section><div class="metric">حجم مصرف‌شده · Used<b>%.2f GiB</b></div><div class="metric">حجم باقی‌مانده · Remaining<b>%s</b></div><div class="metric">حجم کل · Total<b>%s</b></div><div class="metric">انقضا · Expires<b>%s</b></div></section><p>میزان مصرف · Traffic usage</p><progress max="100" value="%.0f"></progress><p>%s</p><footer>لینک اشتراک محرمانه است و نباید عمومی منتشر شود. · Private bearer link.</footer></article><script>const btn=document.getElementById("copy");if(btn)btn.addEventListener("click",async()=>{const v=document.getElementById("subscription");try{await navigator.clipboard.writeText(v.value);btn.textContent="کپی شد · Copied"}catch(e){v.select();btn.textContent="لینک انتخاب شد · Selected"}})</script></body></html>`,name,html.EscapeString(status),used,remainText,quotaText,expiry,ratio,action)
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  if r.Method=="GET" {_,_=w.Write([]byte(body))}
}

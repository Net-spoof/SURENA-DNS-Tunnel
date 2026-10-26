package panel

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

// The existing subscription bearer token protects this customer-facing view.
// No administrator session, extra database table or credential is created.
type subscriberProfile struct {
	Engine  string
	Zone    string
	Link    template.URL // Only after checking the exact expected URI scheme.
	Text    string
	QRPath  string
	Index   int
}

type subscriberPortalPage struct {
	Name             string
	StateFA          string
	StateEN          string
	Active           bool
	Used             string
	Remaining        string
	Total            string
	Expiry           string
	Days             string
	Percent          int
	SubscriptionURL  string
	SubscriptionQR   string
	Profiles         []subscriberProfile
}

func subscriberGiB(n int64) string {
	return fmt.Sprintf("%.2f GiB", float64(n)/(1024*1024*1024))
}

// Validate protocol URI schemes before marking them safe for an HTML link.
// Account names, DNS zones and configuration text remain automatically escaped.
func safeSubscriberURI(engine, raw string) (template.URL, bool) {
	schemes := map[string]string{
		"soar": "soar://",
		"cotten": "cottendns://",
		"masterdns": "masterdns://",
	}
	prefix := schemes[engine]
	if prefix == "" || !strings.HasPrefix(raw, prefix) || len(raw) <= len(prefix) {
		return "", false
	}
	for _, c := range raw {
		if c <= 0x20 || c == 0x7f || c == '<' || c == '>' || c == '"' || c == '\'' || c == '\\' {
			return "", false
		}
	}
	return template.URL(raw), true
}

var subscriberPortalTemplate = template.Must(template.New("subscriber").Parse(`
<!doctype html>
<html lang="fa" dir="rtl">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="color-scheme" content="dark">
<meta name="robots" content="noindex,nofollow,noarchive">
<title>پنل اشتراک · Subscription Portal</title>
<style>
:root{font-family:system-ui,-apple-system,"Segoe UI",Tahoma,sans-serif;color-scheme:dark}
*{box-sizing:border-box}
body{margin:0;min-height:100vh;padding:30px 16px 48px;background:radial-gradient(ellipse at 2% 7%,#1c4563,transparent 39%),radial-gradient(ellipse at 98% 92%,#0c3b42,transparent 37%),#0a1019;color:#eaf5fc}
main{max-width:860px;margin:0 auto}
.masthead{display:flex;justify-content:space-between;align-items:center;gap:12px;margin-bottom:20px}
.masthead .title{font-size:22px;font-weight:850;letter-spacing:-.6px}
.masthead .eyebrow{color:#8daabe;font-size:11px;letter-spacing:.1em;margin-top:4px;direction:ltr;text-align:right}
.status{display:inline-flex;align-items:center;gap:7px;color:#8af3cd;background:#123a32;border:1px solid #367562;border-radius:100px;padding:8px 12px;font-size:12px}
.status.off{background:#443025;color:#ffd2ad;border-color:#80523c}
.status:before{content:"";display:block;width:7px;height:7px;background:currentColor;border-radius:50%}
.panel{padding:28px;background:linear-gradient(135deg,#18293a,#111d2b);border:1px solid #395267;border-radius:24px;box-shadow:0 22px 75px #0008}
.kicker{color:#6ee6cc;font-size:11px;font-weight:700;letter-spacing:.08em;direction:ltr;text-align:right}
h1{font-size:clamp(20px,4vw,29px);margin:10px 0 25px;overflow-wrap:anywhere}
.stats{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:12px}
.stat{min-width:0;border:1px solid #ffffff14;border-radius:15px;background:#ffffff09;padding:15px}
.stat .name{display:block;color:#9cb2c3;font-size:12px;margin-bottom:10px}
.stat strong{display:block;direction:ltr;text-align:right;font-size:16px;overflow-wrap:anywhere}
.usage{margin:24px 0}
.usage-head{display:flex;justify-content:space-between;align-items:center;color:#b7cad9;font-size:12px;margin-bottom:10px}
progress{width:100%;height:11px;appearance:none;-webkit-appearance:none;border:0;border-radius:100px;overflow:hidden;background:#314356}
progress::-webkit-progress-bar{background:#314356;border-radius:100px}
progress::-webkit-progress-value{background:linear-gradient(90deg,#168da7,#5ce0b6)}
progress::-moz-progress-bar{background:#46cfac}
.section{margin-top:28px;border-top:1px solid #ffffff1c;padding-top:23px}
.section-head{display:flex;align-items:center;justify-content:space-between;gap:10px;margin-bottom:17px}
h2{font-size:18px;margin:0}
.subtitle{font-size:12px;color:#96afc0;margin-top:5px;direction:ltr;text-align:right}
.count{border:1px solid #466578;color:#bed6e2;background:#ffffff08;padding:6px 10px;border-radius:999px;font-size:11px;white-space:nowrap}
.actions{display:flex;gap:10px;flex-wrap:wrap;margin:13px 0}
.btn,.button,.qr-toggle{font:inherit;font-size:12px;font-weight:750;display:inline-flex;justify-content:center;align-items:center;gap:5px;border:1px solid #4c667b;background:#223e51;color:#ebf7fe;border-radius:11px;padding:11px 13px;cursor:pointer;text-decoration:none;min-height:43px;text-align:center}
.btn.primary,.button.primary{color:#061f23;background:#5fe2c7;border-color:#5fe2c7}
.btn:hover,.button:hover,.qr-toggle:hover{filter:brightness(1.12)}.phone-icon{font-size:16px;line-height:1}
button:focus-visible,summary:focus-visible,a:focus-visible{outline:2px solid #7bdccf;outline-offset:3px}
.subfield,.configfield{display:block;width:100%;background:#0b1522;border:1px solid #355065;border-radius:11px;color:#d7e6f1;padding:12px;font:12px ui-monospace,SFMono-Regular,Consolas,monospace;direction:ltr;text-align:left;overflow-wrap:anywhere}
.configfield{resize:vertical;min-height:76px}
.note{font-size:12px;color:#9bb1c1;line-height:1.85;margin:12px 0}
.subqr{text-align:center;margin:17px auto 0;padding-top:15px;border-top:1px solid #ffffff13}
.subqr img,.qr-details img{display:block;margin:0 auto;background:white;border-radius:13px;padding:9px;width:178px;max-width:100%;height:auto;image-rendering:auto}
.profile-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:14px}
.profile-card{min-width:0;background:#0c1b2a;border:1px solid #314e63;border-radius:18px;padding:17px}
.profile-top{display:flex;justify-content:space-between;align-items:center;gap:7px;margin-bottom:13px}
.engine{display:inline-flex;padding:7px 11px;border-radius:999px;background:#13443f;color:#79f1d7;font-size:12px;font-weight:800;direction:ltr}
.engine.cotten{background:#26385e;color:#a3c3ff}
.engine.masterdns{background:#41305f;color:#d0afff}
.zone{color:#91adbe;font-size:11px;direction:ltr;text-align:right;overflow-wrap:anywhere;min-width:0}
.profile-card .actions{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));margin:12px 0 0}
.profile-card .actions>*{min-width:0}
.qr-details{margin-top:12px}
.qr-details summary{list-style:none;user-select:none}
.qr-details summary::-webkit-details-marker{display:none}
.qr-details[open] summary{border-color:#70e7cf;color:#89f6dc}
.qr-content{margin-top:12px;padding:17px 10px;background:#101f2e;border:1px solid #345063;border-radius:13px;text-align:center}
.qr-content img{width:205px}
.expand{margin-top:11px}
.privacy{margin-top:20px;text-align:center;font-size:11px;line-height:1.8;color:#91a6b8}
dialog{width:min(92vw,430px);padding:22px;border-radius:21px;background:#162737;border:1px solid #517184;color:#f0faff;box-shadow:0 15px 55px #000c}
dialog::backdrop{background:#000c}
.dialog-title{display:flex;align-items:center;justify-content:space-between;margin-bottom:15px;gap:12px}
dialog img{display:block;width:min(100%,345px);margin:0 auto;background:#fff;border-radius:14px;padding:12px}
.small{font-size:11px;color:#98abb9;text-align:center;margin:12px 0}.import-help{font-size:14px;line-height:2;color:#c7d9e4}.import-help ol{padding-right:20px;margin:11px 0}.import-help li{margin:9px 0}.import-help strong{color:#7be9cb}.dialog-actions{display:flex;gap:10px;flex-wrap:wrap;margin-top:16px}.dialog-actions>*{flex:1}.import-feedback{font-size:12px;line-height:1.8;color:#8edcc7}.import-feedback.warn{color:#ffce9c}
@media(max-width:720px){.stats{grid-template-columns:repeat(2,minmax(0,1fr))}}
@media(max-width:560px){body{padding:16px 11px 34px}.panel{padding:19px}.profile-grid{grid-template-columns:1fr}.masthead .title{font-size:19px}.profile-card{padding:16px}.actions>*{flex:1}.stats{gap:9px}.stat{padding:12px}.stat strong{font-size:14px}}
</style>
</head>
<body>
<main>
<header class="masthead">
<div><div class="title">پنل اشتراک</div><div class="eyebrow">SUBSCRIPTION MANAGEMENT</div></div>
<span class="status {{if not .Active}}off{{end}}">{{.StateFA}} · {{.StateEN}}</span>
</header>
<article class="panel">
<div class="kicker">ACCOUNT OVERVIEW</div><h1>{{.Name}}</h1>
<div class="stats">
<div class="stat"><span class="name">مصرف‌شده · Used</span><strong>{{.Used}}</strong></div>
<div class="stat"><span class="name">باقی‌مانده · Left</span><strong>{{.Remaining}}</strong></div>
<div class="stat"><span class="name">حجم کل · Total</span><strong>{{.Total}}</strong></div>
<div class="stat"><span class="name">انقضا · Expiry</span><strong>{{.Expiry}}</strong></div>
</div>
<div class="usage"><div class="usage-head"><span>میزان مصرف · Usage</span><b>{{.Percent}}%</b></div><progress value="{{.Percent}}" max="100">{{.Percent}}%</progress></div>
{{if .Active}}
<section class="section" aria-labelledby="subscription-heading">
<div class="section-head"><div><h2 id="subscription-heading">لینک اشتراک</h2><div class="subtitle">Subscription tools</div></div></div>
<input class="subfield" id="subscription-url" readonly type="text" value="{{.SubscriptionURL}}" aria-label="Subscription URL">
<div class="actions">
<button class="button primary" id="add-to-app" type="button"><span class="phone-icon" aria-hidden="true">📱</span>افزودن به برنامه</button>
<button class="button copy" type="button" data-copy="subscription-url">کپی لینک اشتراک · Copy subscription</button>
</div>
<details class="qr-details subqr"><summary class="qr-toggle">نمایش QR اشتراک · Show subscription QR</summary>
<div class="qr-content"><img loading="lazy" alt="Subscription QR" src="{{.SubscriptionQR}}" width="205" height="205"><p class="note">برای افزودن اشتراک به اپ، اسکن کنید.</p></div>
</details>
</section>
<section class="section" aria-labelledby="profiles-heading">
<div class="section-head"><div><h2 id="profiles-heading">کانفیگ‌های اتصال</h2><div class="subtitle">Connection profiles</div></div><span class="count">{{len .Profiles}} کانفیگ</span></div>
{{if .Profiles}}
<div class="profile-grid">
{{range .Profiles}}
<article class="profile-card">
<div class="profile-top"><span class="engine {{.Engine}}">{{if eq .Engine "soar"}}Soar{{else if eq .Engine "cotten"}}CottenDNS{{else}}MasterDNS{{end}}</span><span class="zone" dir="ltr">{{.Zone}}</span></div>
<textarea class="configfield" id="config-{{.Index}}" rows="3" readonly spellcheck="false" aria-label="Config {{.Engine}}">{{.Text}}</textarea>
<div class="actions">
<button type="button" class="button primary copy" data-copy="config-{{.Index}}">کپی کانفیگ · Copy</button>
<a class="btn" href="{{.Link}}" rel="noreferrer">بازکردن · Open</a>
</div>
<details class="qr-details">
<summary class="qr-toggle">نمایش QR اختصاصی · Show QR</summary>
<div class="qr-content">
<img loading="lazy" src="{{.QRPath}}" alt="{{.Engine}} config QR" width="205" height="205">
<button type="button" class="button expand qr-expand" data-qr="{{.QRPath}}" data-label="{{.Engine}}">بزرگ‌نمایی QR · Enlarge</button>
<p class="note">این QR فقط کانفیگ همین پروتکل را دارد.</p>
</div>
</details>
</article>
{{end}}
</div>
{{else}}<p class="note">برای این حساب کانفیگ فعالی ثبت نشده است.</p>{{end}}
</section>
{{else}}
<section class="section"><h2>اشتراک غیرفعال است</h2><p class="note">حساب غیرفعال، منقضی یا فاقد حجم است؛ لینک کانفیگ و QR اتصال نمایش داده نمی‌شوند.</p></section>
{{end}}
</article>
<p class="privacy">لینک اشتراک و QRهای اتصال محرمانه‌اند. آن‌ها را در پیام عمومی یا تصویر منتشر نکنید.</p>
</main>
<dialog id="app-dialog" aria-labelledby="app-dialog-title">
<div class="dialog-title"><strong id="app-dialog-title">افزودن اشتراک به برنامه</strong><button id="app-close" type="button" class="button">بستن · Close</button></div>
<div class="import-help"><p class="import-feedback" id="import-feedback" role="status"></p><ol><li>برنامه سازگار را باز کن.</li><li>به قسمت <strong>Profiles → Connection → Import</strong> برو.</li><li>لینک اشتراک را <strong>Paste</strong> کن و واردکردن را تأیید کن.</li></ol><p class="note">ورود مستقیم کانفیگ‌های تکی از طریق دکمه Open در کارت هر پروتکل امکان‌پذیر است. برای اشتراک، لینک باید داخل بخش Import جای‌گذاری شود.</p></div>
<div class="dialog-actions"><button type="button" id="app-copy-again" class="button primary copy" data-copy="subscription-url">کپی دوباره لینک</button><button type="button" id="app-done" class="button">متوجه شدم</button></div>
</dialog>
<dialog id="qr-dialog" aria-label="QR code">
<div class="dialog-title"><strong id="qr-caption"></strong><button id="qr-close" type="button" class="button">بستن · Close</button></div>
<img id="qr-big" alt="Full-size protocol QR code">
<p class="small">این QR فقط برای کانفیگ انتخاب‌شده است.</p>
</dialog>
<script>
(function(){
function setCopied(btn,message){const label=btn.dataset.label||(btn.dataset.label=btn.textContent);btn.textContent=message;setTimeout(()=>{btn.textContent=label},1600)}
async function copy(input,btn){
if(!input)return false;
try{if(!navigator.clipboard||!window.isSecureContext)throw new Error("clipboard unavailable");await navigator.clipboard.writeText(input.value);setCopied(btn,"کپی شد · Copied");return true}
catch(e){input.focus();input.select();let ok=false;try{ok=document.execCommand("copy")}catch(err){ok=false}setCopied(btn,ok?"کپی شد · Copied":"انتخاب شد · Selected");return ok}
}
document.querySelectorAll(".copy").forEach(b=>b.addEventListener("click",()=>copy(document.getElementById(b.dataset.copy),b)));
const add=document.getElementById("add-to-app"),appDialog=document.getElementById("app-dialog");
if(add) add.addEventListener("click",async()=>{
 const ok=await copy(document.getElementById("subscription-url"),add);
 const msg=document.getElementById("import-feedback");
 msg.textContent=ok?"لینک اشتراک کپی شد؛ آن را در بخش Import برنامه جای‌گذاری کن.":"کپی خودکار ممکن نشد؛ از دکمه کپی لینک اشتراک استفاده کن و سپس در برنامه Paste کن.";
 msg.classList.toggle("warn",!ok);
 if(typeof appDialog.showModal==="function") appDialog.showModal(); else appDialog.setAttribute("open","");
});
const closeApp=()=>{if(typeof appDialog.close==="function")appDialog.close();else appDialog.removeAttribute("open")};
document.getElementById("app-close").addEventListener("click",closeApp);
document.getElementById("app-done").addEventListener("click",closeApp);
appDialog.addEventListener("click",e=>{if(e.target===appDialog)closeApp()});
const dialog=document.getElementById("qr-dialog"),image=document.getElementById("qr-big"),caption=document.getElementById("qr-caption");
document.querySelectorAll(".qr-expand").forEach(b=>b.addEventListener("click",()=>{
image.src=b.dataset.qr;caption.textContent=b.dataset.label+" · QR";
if(typeof dialog.showModal==="function")dialog.showModal();else window.open(image.src,"_blank","noopener,noreferrer");
}));
document.getElementById("qr-close").addEventListener("click",()=>dialog.close());
dialog.addEventListener("click",e=>{if(e.target===dialog)dialog.close()});
})();
</script>
</body></html>
`))

// Each QR is generated server-side from the EXACT same profile URL displayed
// in the matching card. The query string never contains the config itself.
func (h *HTTP) serveDNSFabricPortal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := strings.TrimPrefix(r.URL.Path, "/u/")
	account, cfg, err := h.store.Subscription(token)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	state := account.Status(time.Now())
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; object-src 'none'; form-action 'none'; base-uri 'none'; frame-ancestors 'none'")

	qrType := r.URL.Query().Get("qr")
	if qrType != "" {
		if state != "active" {
			http.NotFound(w, r)
			return
		}
		var qrValue string
		switch qrType {
		case "1":
			qrValue = strings.Replace(h.subscriptionURL(r, token), "/sub/", "/s/", 1)
		case "profile":
			engine, zone := r.URL.Query().Get("engine"), r.URL.Query().Get("zone")
			validZone := false
			for _, candidate := range cfg.orderedZones() {
				if candidate == zone {
					validZone = true
					break
				}
			}
			if !validZone {
				http.NotFound(w, r)
				return
			}
			profiles, e := h.store.profilesAt(account, cfg, zone)
			if e != nil {
				http.Error(w, "QR unavailable", http.StatusInternalServerError)
				return
			}
			qrValue = profiles[engine]
			if _, ok := safeSubscriberURI(engine, qrValue); !ok {
				http.NotFound(w, r)
				return
			}
		default:
			http.NotFound(w, r)
			return
		}
		png, e := qrcode.Encode(qrValue, qrcode.Medium, 512)
		if e != nil {
			http.Error(w, "QR unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		if r.Method == http.MethodGet {
			_, _ = w.Write(png)
		}
		return
	}

	labels := map[string][2]string{
		"active": {"فعال", "Active"},
		"disabled": {"غیرفعال", "Disabled"},
		"expired": {"منقضی", "Expired"},
		"exhausted": {"حجم تمام‌شده", "Exhausted"},
	}
	names := labels[state]
	used := account.Used()
	page := subscriberPortalPage{
		Name: account.Name, StateFA: names[0], StateEN: names[1],
		Active: state == "active", Used: subscriberGiB(used),
		Total: "نامحدود", Remaining: "نامحدود", Expiry: "بدون انقضا", Days: "∞",
		SubscriptionURL: strings.Replace(h.subscriptionURL(r, token), "/sub/", "/s/", 1),
		SubscriptionQR: "/u/" + token + "?qr=1",
	}
	if account.QuotaBytes > 0 {
		page.Total = subscriberGiB(account.QuotaBytes)
		page.Remaining = subscriberGiB(max(int64(0), account.QuotaBytes-used))
		page.Percent = int(math.Round(100 * float64(used) / float64(account.QuotaBytes)))
		if page.Percent > 100 {
			page.Percent = 100
		}
	}
	if account.ExpiresAt > 0 {
		expiry := time.UnixMilli(account.ExpiresAt).UTC()
		page.Expiry = expiry.Format("2006-01-02")
		remaining := int(math.Ceil(time.Until(expiry).Hours() / 24))
		if remaining < 0 {
			remaining = 0
		}
		page.Days = fmt.Sprint(remaining)
	}
	if page.Active {
		for _, zone := range cfg.orderedZones() {
			profiles, e := h.store.profilesAt(account, cfg, zone)
			if e != nil {
				http.Error(w, "Profiles unavailable", http.StatusInternalServerError)
				return
			}
			for _, engine := range profileEngines {
				raw := profiles[engine]
				if raw == "" {
					continue
				}
				uri, ok := safeSubscriberURI(engine, raw)
				if !ok {
					http.Error(w, "Unsupported profile", http.StatusInternalServerError)
					return
				}
				page.Profiles = append(page.Profiles, subscriberProfile{
					Engine: engine, Zone: zone, Link: uri, Text: raw,
					Index: len(page.Profiles) + 1,
					QRPath: "/u/" + token + "?qr=profile&engine=" +
						url.QueryEscape(engine) + "&zone=" + url.QueryEscape(zone),
				})
			}
		}
	}
	var body bytes.Buffer
	if err := subscriberPortalTemplate.Execute(&body, page); err != nil {
		http.Error(w, "Portal unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == http.MethodGet {
		_, _ = w.Write(body.Bytes())
	}
}

package panel

import (
 "bytes"
 "net/http/httptest"
 "net/url"
 qrcode "github.com/skip2/go-qrcode"
 "strings"
 "testing"
)

func TestDNSFabricSubscriptionPortal(t *testing.T) {
 _, store, h := testBrowser(t)
 c := account(t, store, 1073741824)
 token, err := store.SubscriptionToken(c.ID,false)
 if err!=nil {t.Fatal(err)}
 legacy := httptest.NewRecorder()
 h.ServeHTTP(legacy,httptest.NewRequest("GET","http://example.test/sub/"+token,nil))
 short := httptest.NewRecorder()
 h.ServeHTTP(short,httptest.NewRequest("GET","http://example.test/s/"+token,nil))
 if legacy.Code!=200 || short.Code!=200 || legacy.Body.String()!=short.Body.String() {
   t.Fatal("short subscription changed legacy output")
 }
 if short.Header().Get("Subscription-Userinfo") != legacy.Header().Get("Subscription-Userinfo") {
   t.Fatal("app short link changed usage metadata")
 }
 // Browser Accept headers display HTML dashboard at the SAME /s/ URL.
 browser := httptest.NewRequest("GET","http://example.test/s/"+token,nil)
 browser.Header.Set("Accept","text/html,application/xhtml+xml,*/*")
 browser.Header.Set("Sec-Fetch-Dest","document")
 webpage:=httptest.NewRecorder()
 h.ServeHTTP(webpage,browser)
 if webpage.Code!=200 || !strings.Contains(webpage.Header().Get("Content-Type"),"text/html") || !strings.Contains(webpage.Body.String(),"پنل اشتراک") {
   t.Fatal("browser /s/ navigation must render subscriber portal")
 }
 if !strings.Contains(webpage.Body.String(),"soar://") {t.Fatal("active subscriber portal must include direct profile link")}
 // Explicit format always wins, even in a browser.
 rawRequest := httptest.NewRequest("GET","http://example.test/s/"+token+"?format=base64",nil)
 rawRequest.Header.Set("Accept","text/html")
 forcedRaw:=httptest.NewRecorder()
 h.ServeHTTP(forcedRaw,rawRequest)
 if forcedRaw.Body.String()!=legacy.Body.String() {t.Fatal("explicit raw format changed payload")}
 legacyBrowser:=httptest.NewRequest("GET","http://example.test/sub/"+token,nil)
 legacyBrowser.Header.Set("Accept","text/html")
 legacyRec:=httptest.NewRecorder()
 h.ServeHTTP(legacyRec,legacyBrowser)
 if legacyRec.Body.String()!=legacy.Body.String() {t.Fatal("old subscription now HTML")}
 // Client Accept negotiation must leave /s/ in raw mode.
 appRequest:=httptest.NewRequest("GET","http://example.test/s/"+token,nil)
 appRequest.Header.Set("Accept","*/*")
 appRec:=httptest.NewRecorder()
 h.ServeHTTP(appRec,appRequest)
 if appRec.Body.String()!=legacy.Body.String(){t.Fatal("app client incompatible")}
 p:=httptest.NewRecorder()
 h.ServeHTTP(p,httptest.NewRequest("GET","http://example.test/u/"+token,nil))
 if p.Code!=200 || !strings.Contains(p.Body.String(),"پنل اشتراک") {
   t.Fatal("subscriber portal failed")
 }
 if !strings.Contains(p.Header().Get("X-Robots-Tag"),"noindex") || !strings.Contains(p.Header().Get("X-Robots-Tag"),"nofollow") {t.Fatal("portal must not be indexed")}
 bad:=httptest.NewRecorder()
 h.ServeHTTP(bad,httptest.NewRequest("GET","http://example.test/u/bad",nil))
 if bad.Code!=404 {t.Fatal("invalid bearer token accepted")}
}

func TestSubscriberPageListsAllEnabledProtocolConfigs(t *testing.T) {
 _, store, h := testBrowser(t)
 c, err := store.Create(ClientInput{Name:"Multi protocol test", Engines:"all", Enabled:true, QuotaBytes:1073741824, StreamLimit:16})
 if err != nil { t.Fatal(err) }
 credential, err := store.SubscriptionToken(c.ID,false)
 if err != nil { t.Fatal(err) }
 page := httptest.NewRecorder()
 req := httptest.NewRequest("GET","http://example.test/s/"+credential,nil)
 req.Header.Set("Accept","text/html")
 req.Header.Set("Sec-Fetch-Dest","document")
 h.ServeHTTP(page,req)
 if page.Code != 200 { t.Fatalf("unexpected page HTTP %d",page.Code) }
 body := page.Body.String()
 if strings.Contains(body, "DNSFabric") || strings.Contains(body, "WhiteDNS") || strings.Contains(body, "WHITE DNS") {
   t.Fatal("subscriber-facing branding must be generic")
 }
 client, cfg, err := store.Subscription(credential)
 if err != nil { t.Fatal(err) }
 count:=0
 for _, zone := range cfg.orderedZones() {
   profiles, err := store.profilesAt(client,cfg,zone)
   if err != nil { t.Fatal(err) }
   for _, engine := range profileEngines {
     link:=profiles[engine]
     if link=="" {continue}
     count++
     if !strings.Contains(body, link) {t.Fatalf("missing raw %s config from subscriber view",engine)}
     qrPath:="/u/"+credential+"?qr=profile&engine="+url.QueryEscape(engine)+"&zone="+url.QueryEscape(zone)
     qr:=httptest.NewRecorder()
     h.ServeHTTP(qr,httptest.NewRequest("GET","http://example.test"+qrPath,nil))
     if qr.Code!=200 || qr.Header().Get("Content-Type")!="image/png" {t.Fatalf("missing QR for %s: %d",engine,qr.Code)}
     expectedQR,e:=qrcode.Encode(link,qrcode.Medium,512)
     if e!=nil {t.Fatal(e)}
     if !bytes.Equal(qr.Body.Bytes(),expectedQR) {t.Fatalf("QR payload mismatch for %s; must encode direct config",engine)}
   }
 }
 if count != 3 {t.Fatalf("expected 3 config profiles, got %d",count)}
 if strings.Count(body,"class=\"button copy\" data-copy=\"config-") != count {
   t.Fatal("each config card must have its own copy button")
 }
 if strings.Count(body,"class=\"qr-details\"") != count {
   t.Fatal("each config card must have independently expandable QR")
 }
 if strings.Count(body,"class=\"button expand qr-expand\"") != count {
   t.Fatal("each config must support enlarging its QR")
 }
 if strings.Count(body,"📱 افزودن به برنامه") != count {t.Fatal("every profile must have an Add to App direct deep link")}
 if strings.Count(body,"class=\"btn primary profile-open\"") != count {t.Fatal("every profile must have a direct deep link button")}
 if strings.Contains(body,"#ZgotmplZ") {t.Fatal("Go HTML templates rejected legitimate protocol deep link")}
 if !strings.Contains(body,"Copy subscription") {t.Fatal("subscription copy button should be separately named")}
 if !strings.Contains(body,"id=\"add-to-app\"") || !strings.Contains(body,"کپی و راهنمای افزودن اشتراک") {
   t.Fatal("subscription control must be clearly distinguished from direct profile Add to App")
 }
 if !strings.Contains(body,"id=\"app-dialog\"") || !strings.Contains(body,"Profiles → Connection → Import") ||
   !strings.Contains(body,"id=\"app-copy-again\"") {
   t.Fatal("subscription import must display accurate copy-and-paste instructions")
 }
 if strings.Contains(body,"دریافت اشتراک · Import</a>") || strings.Contains(body,"بازکردن · Open</a>") {
   t.Fatal("misleading legacy action labels must not remain")
 }
 if !strings.Contains(body,"id=\"subscription-url\"") || !strings.Contains(body,"data-copy=\"subscription-url\"") {
   t.Fatal("separate subscription copy control must remain")
 }
 // A valid bearer token is not enough to retrieve arbitrary engines or DNS zones.
 for _, query := range []string{"?qr=profile&engine=invalid&zone=example.com","?qr=profile&engine=soar&zone=attacker.example.com","?qr=invalid"} {
   invalid:=httptest.NewRecorder()
   h.ServeHTTP(invalid,httptest.NewRequest("GET","http://example.test/u/"+credential+query,nil))
   if invalid.Code!=404 {t.Fatalf("unexpected QR endpoint response for invalid request: %d",invalid.Code)}
 }
 // Disabled users retain a status page, but not connection secrets.
 clientInput:=ClientInput{Name:c.Name,Engines:"all",Enabled:false,QuotaBytes:c.QuotaBytes,StreamLimit:c.StreamLimit,Revision:c.Revision}
 if _,err := store.Edit(c.ID,clientInput);err!=nil {t.Fatal(err)}
 after:=httptest.NewRecorder()
 h.ServeHTTP(after,httptest.NewRequest("GET","http://example.test/u/"+credential,nil))
 if after.Code!=200 || strings.Contains(after.Body.String(),"soar://") ||
   strings.Contains(after.Body.String(),"cottendns://") || strings.Contains(after.Body.String(),"masterdns://") {
   t.Fatal("disabled account should not display configuration secrets")
 }
 if strings.Contains(after.Body.String(),"id=\"add-to-app\"") {t.Fatal("inactive accounts cannot add subscription to app")}
 blockedQR:=httptest.NewRecorder()
 h.ServeHTTP(blockedQR,httptest.NewRequest("GET","http://example.test/u/"+credential+"?qr=profile&engine=soar&zone="+url.QueryEscape(cfg.orderedZones()[0]),nil))
 if blockedQR.Code!=404 {t.Fatal("disabled account must not expose QR credentials")}
}

func TestOfficialClientDownloadActionsAreSafe(t *testing.T) {
 _, store, h := testBrowser(t)
 account := account(t, store, 1073741824)
 token, err := store.SubscriptionToken(account.ID,false)
 if err != nil { t.Fatal(err) }
 expectedAPK := "https://github.com/WhiteDNS/WhiteDNS-Android/releases/download/1.6.14/WhiteDNS-1.6.14-universal.apk"
 for _, uri := range []string{"/u/"+token, "/s/"+token} {
   request:=httptest.NewRequest("GET","http://example.test"+uri,nil)
   if strings.HasPrefix(uri,"/s/") {request.Header.Set("Accept","text/html");request.Header.Set("Sec-Fetch-Dest","document")}
   page:=httptest.NewRecorder()
   h.ServeHTTP(page,request)
   if page.Code!=200 {t.Fatalf("%s unexpected HTTP %d",uri,page.Code)}
   html:=page.Body.String()
   if !strings.Contains(html,expectedAPK) || !strings.Contains(html,`id="android-app-download"`) {
     t.Fatalf("%s must offer official direct-download universal Android APK",uri)
   }
   if !strings.Contains(html,"https://github.com/WhiteDNS/WhiteDNS-Android/releases/latest") {
     t.Fatal("missing official latest version fallback")
   }
   if !strings.Contains(html,"آیفون / iOS") || !strings.Contains(html,"نسخه iOS فعلاً موجود نیست") {
     t.Fatal("iOS availability must be accurately labeled")
   }
   if strings.Contains(html,"apps.apple.com/") || strings.Contains(html,"testflight.apple.com/") {
     t.Fatal("unsupported iOS download URL must not be advertised")
   }
   if !strings.Contains(html,`class="app-unavailable" aria-disabled="true"`) {
     t.Fatal("iOS unavailable indicator should not be a clickable download")
   }
   if strings.Contains(html,`href="https://github.com/WhiteDNS/WhiteDNS-Android/releases/download/1.6.14/WhiteDNS-1.6.14-universal.apk/u/`) {
     t.Fatal("public download must not include bearer credentials")
   }
 }
 // Inactive users should still be able to download the app, without
 // disclosing any active protocol configuration.
 changed,err:=store.Edit(account.ID,ClientInput{
   Name:account.Name,Engines:account.Engines,Enabled:false,
   QuotaBytes:account.QuotaBytes,StreamLimit:account.StreamLimit,
   Revision:account.Revision,
 })
 if err != nil {t.Fatal(err)}
 _ = changed
 disabled:=httptest.NewRecorder()
 h.ServeHTTP(disabled,httptest.NewRequest("GET","http://example.test/u/"+token,nil))
 if disabled.Code!=200 || !strings.Contains(disabled.Body.String(),expectedAPK) {
   t.Fatal("disabled subscriber page should still offer client download")
 }
 if strings.Contains(disabled.Body.String(),`class="profile-card"`) {
   t.Fatal("disabled subscriber must not receive any live profiles")
 }
}

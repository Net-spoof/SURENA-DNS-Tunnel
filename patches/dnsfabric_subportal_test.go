package panel

import (
 "net/http/httptest"
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
 if p.Header().Get("X-Robots-Tag")!="noindex, nofollow" {t.Fatal("portal must not be indexed")}
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
   }
 }
 if count != 3 {t.Fatalf("expected 3 config profiles, got %d",count)}
 if strings.Count(body,"class='copy-config'") != count {t.Fatal("every config must have its own copy button")}
 if !strings.Contains(body,"Copy subscription") {t.Fatal("subscription copy button should be separately named")}
 // Disabled users retain a status page, but not connection secrets.
 clientInput:=ClientInput{Name:c.Name,Engines:"all",Enabled:false,QuotaBytes:c.QuotaBytes,StreamLimit:c.StreamLimit,Revision:c.Revision}
 if _,err := store.Edit(c.ID,clientInput);err!=nil {t.Fatal(err)}
 after:=httptest.NewRecorder()
 h.ServeHTTP(after,httptest.NewRequest("GET","http://example.test/u/"+credential,nil))
 if after.Code!=200 || strings.Contains(after.Body.String(),"soar://") ||
   strings.Contains(after.Body.String(),"cottendns://") || strings.Contains(after.Body.String(),"masterdns://") {
   t.Fatal("disabled account should not display configuration secrets")
 }
}

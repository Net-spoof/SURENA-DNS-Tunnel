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
 p:=httptest.NewRecorder()
 h.ServeHTTP(p,httptest.NewRequest("GET","http://example.test/u/"+token,nil))
 if p.Code!=200 || !strings.Contains(p.Body.String(),"DNSFabric") {
   t.Fatal("subscriber portal failed")
 }
 if p.Header().Get("X-Robots-Tag")!="noindex, nofollow" {t.Fatal("portal must not be indexed")}
 bad:=httptest.NewRecorder()
 h.ServeHTTP(bad,httptest.NewRequest("GET","http://example.test/u/bad",nil))
 if bad.Code!=404 {t.Fatal("invalid bearer token accepted")}
}

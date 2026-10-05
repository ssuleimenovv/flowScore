package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func corsServer(t *testing.T) *httptest.Server {
	t.Helper()
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	srv := httptest.NewServer(CORS(ok, []string{"flowscore.vercel.app", "flowscore-*.vercel.app"}))
	t.Cleanup(srv.Close)
	return srv
}

func send(t *testing.T, method, url, origin string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", origin)
	if method == http.MethodOptions {
		req.Header.Set("Access-Control-Request-Method", "POST")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	return res
}

func TestCORSAllowsTheFrontend(t *testing.T) {
	srv := corsServer(t)

	for _, origin := range []string{"https://flowscore.vercel.app", "https://flowscore-git-deploy-me.vercel.app"} {
		res := send(t, http.MethodGet, srv.URL, origin)
		if got := res.Header.Get("Access-Control-Allow-Origin"); got != origin {
			t.Errorf("%s: Allow-Origin = %q", origin, got)
		}
	}
}

func TestCORSAnswersPreflight(t *testing.T) {
	srv := corsServer(t)

	res := send(t, http.MethodOptions, srv.URL, "https://flowscore.vercel.app")
	if res.StatusCode != http.StatusNoContent || res.Header.Get("Access-Control-Allow-Headers") != "Content-Type" {
		t.Errorf("preflight: status %d, headers %v", res.StatusCode, res.Header)
	}
}

func TestCORSIgnoresOtherSites(t *testing.T) {
	srv := corsServer(t)

	res := send(t, http.MethodGet, srv.URL, "https://evil.vercel.app")
	if got := res.Header.Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Allow-Origin = %q for a site not on the list", got)
	}
}

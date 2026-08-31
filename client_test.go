package nsl

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientCreateAndDeleteUseV2Contract(t *testing.T) {
	t.Parallel()
	var deleteIfMatch string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/apps":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"app-id","node_id":"node-id","generation":1,"name":"api","target_url":"http://api:4000","public_url":"https://api--mac.example.com","routes":[],"enabled":true,"created_at":"2026-08-28T00:00:00Z","updated_at":"2026-08-28T00:00:00Z"}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v2/apps/app-id":
			deleteIfMatch = r.Header.Get("If-Match")
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	client := NewClient(server.URL)
	app, err := client.Create(AppInput{Name: "api", TargetURL: "http://api:4000"})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Delete(app); err != nil {
		t.Fatal(err)
	}
	if deleteIfMatch != `"app:app-id:1"` {
		t.Fatalf("If-Match = %q", deleteIfMatch)
	}
}

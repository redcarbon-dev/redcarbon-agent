package routines

import (
	"net/http"
	"reflect"
	"testing"
)

func TestHeaderNamesOmitsValues(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Basic c2VjcmV0")
	h.Set("X-Api-Key", "secret")
	h.Set("Content-Type", "application/json")

	got := headerNames(h)

	want := []string{"Authorization", "Content-Type", "X-Api-Key"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("headerNames() = %v, want %v", got, want)
	}
}

func TestRedactedURLMasksSecrets(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"no query", "https://wazuh:9200/wazuh-alerts-*/_search", "https://wazuh:9200/wazuh-alerts-*/_search"},
		{"query values", "https://vendor.example/api?apikey=secret&page=2", "https://vendor.example/api?apikey=xxxxx&page=xxxxx"},
		{"userinfo password", "https://user:secret@vendor.example/api", "https://user:xxxxx@vendor.example/api"},
		{"unparseable", "http://[::1", "<unparseable URL>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactedURL(tt.raw); got != tt.want {
				t.Fatalf("redactedURL(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

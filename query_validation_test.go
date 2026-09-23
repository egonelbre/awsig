package awsig

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVerifiersRejectMalformedQuery(t *testing.T) {
	for _, v := range []struct {
		name   string
		verify func(*http.Request) error
	}{
		{"V2", func(r *http.Request) error { _, err := NewV2[struct{}](nil).Verify(r, ""); return err }},
		{"V4", func(r *http.Request) error { _, err := NewV4[struct{}](nil, V4Config{}).Verify(r); return err }},
		{"V2V4", func(r *http.Request) error { _, err := NewV2V4[struct{}](nil, V4Config{}).Verify(r, ""); return err }},
	} {
		for _, raw := range []string{"bad=%ZZ", "bad%=x", "a=1;b=2", "ok=1&bad=%", "X-Amz-Algorithm=AWS4-HMAC-SHA256&bad=%ZZ", "AWSAccessKeyId=key&bad=%ZZ"} {
			for _, mode := range []string{"header", "query", "POST"} {
				t.Run(v.name+"/"+raw+"/"+mode, func(t *testing.T) {
					r := httptest.NewRequest("GET", "https://example.com/?"+raw, strings.NewReader("unread body"))
					if mode == "header" {
						r.Header.Set("Authorization", "invalid")
					}
					if mode == "POST" {
						r.Method = "POST"
						r.Header.Set("Content-Type", "multipart/form-data; boundary=test")
					}
					if err := v.verify(r); !errors.Is(err, ErrInvalidRequest) {
						t.Fatalf("got %v, want ErrInvalidRequest", err)
					}
				})
			}
		}
	}
}

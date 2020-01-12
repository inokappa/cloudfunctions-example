package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsJson(t *testing.T) {
	tests := []struct {
		body string
		want bool
	}{
		{body: `[{"email": "foo@example.com"}]`, want: true},
		{body: `[{"email": "foo@example.com"},{"email": "bar@example.com"}]`, want: true},
		{body: "foo", want: false},
	}
	for _, test := range tests {
		actual := IsJson(test.body)
		expected := test.want
		if actual != expected {
			t.Errorf("got: %v\nwant: %v", actual, expected)
		}
	}
}

func TestGcfWebhookHandler(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{body: `[{"email": "foo@example.com"}]`, want: "ok"},
		{body: `[{"email": "foo@example.com"},{"email": "bar@example.com"}]`, want: "ok"},
		{body: "foo", want: "error"},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/", strings.NewReader(test.body))
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		// writeStorage をテスト用に置き換える
		writeStorage = func(r *http.Request, strBody string) error {
			if !IsJson(strBody) {
				return errors.New("Request Body is not JSON")
			} else {
				return nil
			}
		}
		GcfWebhookHandler(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("GcfWebhookHandler(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}

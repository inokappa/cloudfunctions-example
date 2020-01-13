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
		{body: `[{"name": "foo"}]`, want: true},
		{body: `[{"name": "foo"},{"name": "bar"}]`, want: true},
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

func TestGcfWebhookHandlerString(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{body: `[{"name": "foo"}]`, want: "ok"},
		{body: `[{"name": "foo"},{"name": "bar"}]`, want: "ok"},
		{body: `{"name": "foo"},{"name": "bar"}`, want: "error\n"},
		{body: "foo", want: "error\n"},
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
			t.Errorf("GcfWebhookHandler Response String(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}

func TestGcfWebhookHandlerCode(t *testing.T) {
	tests := []struct {
		body string
		want int
	}{
		{body: `[{"name": "foo"}]`, want: 200},
		{body: `[{"name": "foo"},{"name": "bar"}]`, want: 200},
		{body: `{"name": "foo"},{"name": "bar"}`, want: 400},
		{body: "foo", want: 400},
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

		if got := rr.Code; got != test.want {
			t.Errorf("GcfWebhookHandler Response Code(%q) = %d, want %d", test.body, got, test.want)
		}
	}
}

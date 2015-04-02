package main

import (
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_githubHandler_should_fail_if_not_post(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:9393/_github", nil)
	if err != nil {
		t.Fatal("err = %v, want nil", err)
	}
	w := httptest.NewRecorder()
	config := &Config{}
	githubHandler(w, req, config)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("w.Code = %v, want %v", w.Code, http.StatusMethodNotAllowed)
	}
}

func Test_githubHandler_should_fail_if_not_correct_user_agent(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:9393/_github", nil)
	if err != nil {
		t.Fatal("err = %v, want nil", err)
	}

	w := httptest.NewRecorder()
	config := &Config{}
	githubHandler(w, req, config)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("w.Code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func newRequest(r io.Reader, signature string) (*http.Request, error) {
	req, err := http.NewRequest("POST", "http://localhost:9393/_github", r)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", githubUserAgent+"1234")
	req.Header.Add(githubSignature, signature)

	return req, nil
}

func Test_githubHandler_should_fail_if_hmac_signature_is_invalid(t *testing.T) {
	r := strings.NewReader(validPingRequest)

	req, err := newRequest(r, "junk")
	if err != nil {
		t.Fatal("err = %v, want nil", err)
	}

	w := httptest.NewRecorder()
	config := &Config{}

	githubHandler(w, req, config)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("w.Code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func Test_githubHandler_should_fail_if_hmac_signature_is_invalid_hex_encoding(t *testing.T) {
	r := strings.NewReader(validPingRequest)

	req, err := newRequest(r, "sha1=abc12")
	if err != nil {
		t.Fatal("err = %v, want nil", err)
	}

	w := httptest.NewRecorder()
	config := &Config{}

	githubHandler(w, req, config)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("w.Code = %v, want %v", w.Code, http.StatusBadRequest)
	}

	expectedBody := "Invalid signature."
	if !strings.HasPrefix(w.Body.String(), expectedBody) {
		t.Fatalf("w.Body = '%v', want %v", w.Body, expectedBody)
	}
}

func Test_githubHandler_should_fail_if_hmac_signature_is_signed_incorrectly(t *testing.T) {
	r := strings.NewReader(validPingRequest)
	sig := hex.EncodeToString(sign([]byte(validPingRequest), "123abc"))

	req, err := newRequest(r, "sha1="+sig)
	if err != nil {
		t.Fatal("err = %v, want nil", err)
	}

	w := httptest.NewRecorder()
	config := &Config{
		Github: &Github{
			HookSecret: "abc123",
		},
	}

	githubHandler(w, req, config)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("w.Code = %v, want %v", w.Code, http.StatusBadRequest)
	}

	expectedBody := "Invalid signature."
	if !strings.HasPrefix(w.Body.String(), expectedBody) {
		t.Fatalf("w.Body = '%v', want %v", w.Body, expectedBody)
	}
}

func Test_githubHandler_should_fail_if_event_type_absent(t *testing.T) {
	r := strings.NewReader(validPingRequest)
	sig := hex.EncodeToString(sign([]byte(validPingRequest), "abc123"))

	req, err := newRequest(r, "sha1="+sig)
	if err != nil {
		t.Fatal("err = %v, want nil", err)
	}

	w := httptest.NewRecorder()
	config := &Config{
		Github: &Github{
			HookSecret: "abc123",
		},
	}

	githubHandler(w, req, config)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("w.Code = %v, want %v", w.Code, http.StatusBadRequest)
	}

	expectedBody := "Invalid event type specified."
	if !strings.HasPrefix(w.Body.String(), expectedBody) {
		t.Fatalf("w.Body = '%v', want %v", w.Body, expectedBody)
	}
}

func Test_githubHandler_should_fail_if_event_type_invalid(t *testing.T) {
	r := strings.NewReader(validPingRequest)
	sig := hex.EncodeToString(sign([]byte(validPingRequest), "abc123"))

	req, err := newRequest(r, "sha1="+sig)
	if err != nil {
		t.Fatal("err = %v, want nil", err)
	}
	req.Header.Add(githubEventType, "pong")

	w := httptest.NewRecorder()
	config := &Config{
		Github: &Github{
			HookSecret: "abc123",
		},
	}

	githubHandler(w, req, config)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("w.Code = %v, want %v", w.Code, http.StatusBadRequest)
	}

	expectedBody := "Invalid event type specified."
	if !strings.HasPrefix(w.Body.String(), expectedBody) {
		t.Fatalf("w.Body = '%v', want %v", w.Body, expectedBody)
	}
}

func Test_githubHandler_should_succeed_with_valid_ping(t *testing.T) {
	r := strings.NewReader(validPingRequest)
	sig := hex.EncodeToString(sign([]byte(validPingRequest), "abc123"))

	req, err := newRequest(r, "sha1="+sig)
	if err != nil {
		t.Fatal("err = %v, want nil", err)
	}
	req.Header.Add(githubEventType, "ping")

	w := httptest.NewRecorder()
	config := &Config{
		Github: &Github{
			HookSecret: "abc123",
		},
	}

	githubHandler(w, req, config)

	if w.Code != http.StatusOK {
		t.Fatalf("w.Code = %v, want %v", w.Code, http.StatusOK)
	}

	expectedBody := "OK: 1"
	if !strings.HasPrefix(w.Body.String(), expectedBody) {
		t.Fatalf("w.Body = '%v', want %v", w.Body, expectedBody)
	}
}

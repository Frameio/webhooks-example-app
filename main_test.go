package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func dummyEvent() (*bytes.Buffer, error) {
	event := &Event{
		ID:   "123",
		Name: "ping",
	}

	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)
	err := encoder.Encode(event)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func testReq(timestamp int64) (*http.Request, *httptest.ResponseRecorder) {
	body, err := dummyEvent()
	if err != nil {
		return nil, nil
	}

	r := httptest.NewRequest(http.MethodPost, "http://example.com/ping", body)

	signature, err := computeSignature(r, timestamp, secretKey)
	if err != nil {
		return nil, nil
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Frameio-Request-Timestamp", strconv.FormatInt(timestamp, 10))
	r.Header.Set("X-Frameio-Signature", signature)
	w := httptest.NewRecorder()

	return r, w
}

func TestVerifyTimeExpired(t *testing.T) {
	ts := time.Now().Add(-6 * time.Minute).Unix()
	r, w := testReq(ts)

	handler(w, r)
	res := w.Result()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status code %d, instead got %d.", http.StatusBadRequest, res.StatusCode)
	}
}

func TestVerifyTimeValid(t *testing.T) {
	ts := time.Now().Unix()
	r, w := testReq(ts)

	handler(w, r)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, instead got %d.", http.StatusOK, res.StatusCode)
	}
}

func TestVerifySignatureInvalid(t *testing.T) {
	ts := time.Now().Unix()
	r, w := testReq(ts)

	r.Header.Set("X-Frameio-Signature", "123abc")

	handler(w, r)
	res := w.Result()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status code %d, instead got %d.", http.StatusUnauthorized, res.StatusCode)
	}
}

func TestVerifySignatureValid(t *testing.T) {
	ts := time.Now().Unix()
	r, w := testReq(ts)

	handler(w, r)
	res := w.Result()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, instead got %d.", http.StatusUnauthorized, res.StatusCode)
	}
}

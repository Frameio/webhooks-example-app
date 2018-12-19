package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"
)

var (
	version   = "v0"
	secretKey = os.Getenv("FRAMEIO_SECRET_KEY")
)

// Event represents a payload delivered by a Webhook service.
type Event struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Resource *Resource `json:"resource"`
	Team     *Team     `json:"team,omitempty"`
	Project  *Project  `json:"project,omitempty"`
	User     *User     `json:"user,omitempty"`
}

// Resource represents the affected resource.
type Resource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Team represents the team the webhook is associated with.
type Team struct {
	*Resource
}

// Project represents the project the webhook is associated with.
type Project struct {
	*Resource
}

// User represents the user who triggered the event.
type User struct {
	*Resource
}

func handler(w http.ResponseWriter, r *http.Request) {
	out, err := httputil.DumpRequest(r, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println(string(out))

	// Verify the message has been delivered in the last 5 minutes.
	timestampStr := r.Header.Get("X-Frameio-Request-Timestamp")
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Since(time.Unix(timestamp, 0)) > 5*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify request signature.
	expected := r.Header.Get("X-Frameio-Signature")
	signature, _ := computeSignature(r, timestamp, secretKey)
	if expected != signature {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var event *Event
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Handle webhook here.
	log.Println(event.ID)

	w.WriteHeader(http.StatusOK)
}

// The request includes headers to enable the recipient to validate
// that the request is from Frame.io and that it's been delivered within
// the expected time range. To learn more about how this works, take a
// look at our docs https://docs.frame.io/docs/webhooks#section-security.
func computeSignature(r *http.Request, timestamp int64, secret string) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	copy := body[:]
	r.Body = ioutil.NopCloser(bytes.NewReader(copy))

	msg := fmt.Sprintf("%s:%d:%s", version, timestamp, string(body))

	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(msg))

	result := fmt.Sprintf("%s=%s", version, hex.EncodeToString(h.Sum(nil)))

	return result, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), nil))
}

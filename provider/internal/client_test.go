package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gotify/go-api-client/v2/client/application"
)

const (
	testHost = "my.coolapp.local"
	// the httptest.NewServer always listens on this address.
	localHost = "127.0.0.1"
)

func TestClientHostOverwrite(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host != testHost {
			t.Errorf("Expected \"Host\" to be %q, got %q", testHost, req.Host)
		}

		// Return valid response for test endpoint
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "[]")
	}))
	defer ts.Close()

	host := testHost
	gotify, err := NewAuthedClient(ts.URL, "test", "test", &host)
	if err != nil {
		t.Fatalf("Could not construct client: %v", err.Error())
	}

	params := application.NewGetAppsParams()
	_, err = gotify.Client.Application.GetApps(params, gotify.Auth)
	if err != nil {
		t.Fatalf("Error during test request: %v", err.Error())
	}
}

func TestClientHostDefault(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.Host, localHost) {
			t.Errorf("Expected \"Host\" to start with %q, got %q", localHost, req.Host)
		}

		// Return valid response for test endpoint
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "[]")
	}))
	defer ts.Close()

	gotify, err := NewAuthedClient(ts.URL, "test", "test", nil)
	if err != nil {
		t.Fatalf("Could not construct client: %v", err.Error())
	}

	params := application.NewGetAppsParams()
	_, err = gotify.Client.Application.GetApps(params, gotify.Auth)
	if err != nil {
		t.Fatalf("Error during test request: %v", err.Error())
	}
}

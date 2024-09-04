package internal

import (
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/gotify"
)

type AuthedGotifyClient struct {
	Client *client.GotifyREST
	Auth   runtime.ClientAuthInfoWriter
}

type OverwriteHostTransport struct {
	Host string
	Next http.RoundTripper
}

func (hot *OverwriteHostTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Host = hot.Host
	return hot.Next.RoundTrip(req)
}

func wrapWithHost(host string, wrap http.RoundTripper) http.RoundTripper {
	return &OverwriteHostTransport{Host: host, Next: wrap}
}

func NewAuthedClient(endpoint string, username string, password string, host *string) (*AuthedGotifyClient, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	transport := http.DefaultTransport
	if host != nil {
		transport = wrapWithHost(*host, transport)
	}

	client := gotify.NewClient(url, &http.Client{Transport: transport})
	auth := auth.BasicAuth(username, password)

	return &AuthedGotifyClient{Client: client, Auth: auth}, nil
}

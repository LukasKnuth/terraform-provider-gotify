package main

import (
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/gotify"
)

type AuthedGotifyClient struct {
	client *client.GotifyREST
	auth   runtime.ClientAuthInfoWriter
}

func NewAuthedClient(endpoint string, username string, password string) (*AuthedGotifyClient, error) {
	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	client := gotify.NewClient(url, &http.Client{})
	auth := auth.BasicAuth(username, password)

	return &AuthedGotifyClient{client: client, auth: auth}, nil
}

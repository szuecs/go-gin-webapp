package client

import (
	"testing"

	konfig "github.com/szuecs/go-gin-webapp/conf/client"
)

func GetTestClient() *Client {
	cfg := &konfig.ClientConfig{
		URL:           "http://localhost:8080",
		Debug:         false,
		Oauth2Enabled: false,
		Username:      "test-user",
	}
	return &Client{
		Config:      cfg,
		AccessToken: "",
	}
}

func TestGetUsername(t *testing.T) {
	cli := GetTestClient()
	userName := cli.GetUsername("")
	if userName != cli.Config.Username {
		t.Fatal("If defined Username in config and no overwrite was passed it should return the configured username.")
	}
	overwriteUser := "foo"
	userName = cli.GetUsername(overwriteUser)
	if cli.Config.Username != overwriteUser {
		t.Fatal("If overwrite was passed it should overwrite the config.")
	}
	if userName != overwriteUser {
		t.Fatal("If overwrite was passed it return the same.")
	}
}

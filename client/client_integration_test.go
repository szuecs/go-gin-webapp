package client

import (
	"net/url"
	"testing"
)

func TestGetRoot(t *testing.T) {
	cli := GetTestClient()
	u, err := url.Parse(cli.Config.URL)
	if err != nil {
		t.Fatal("Could not parse URL provided for the test client:", cli.Config.URL)
	}
	u.Path = "/"
	cli.Get(u)
}

func TestGetHealth(t *testing.T) {
	cli := GetTestClient()
	u, err := url.Parse(cli.Config.URL)
	if err != nil {
		t.Fatal("Could not parse URL provided for the test client:", cli.Config.URL)
	}
	u.Path = "/health"
	cli.Get(u)
}

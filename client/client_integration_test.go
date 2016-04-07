package client

import (
	"fmt"
	"testing"
	"time"

	konfig "github.com/zalando-techmonkeys/chimp/conf/client"
	. "github.com/zalando-techmonkeys/chimp/types"
)

var wait bool

func getClient() *Client {
	clusters := map[string]*konfig.Cluster{}
	clusters["LOCAL"] = &konfig.Cluster{IP: "localhost", Port: 8082}
	clientino := Client{Config: &konfig.ClientConfig{Clusters: clusters}, Scheme: "http"}
	return &clientino
}

func TestGetList(t *testing.T) {
	clientino := getClient()
	clientino.ListDeploy(true)
}

func TestGetInfoExisting(t *testing.T) {
	clientino := getClient()
	clientino.InfoDeploy("p5", false) //TODO: refactor to auto-initialize config to test with real enviroment
}

func TestGetInfoNonExisting(t *testing.T) {
	clientino := getClient()
	clientino.InfoDeploy("nonexistingservice", false)
}

func TestDeploy(t *testing.T) {
	clientino := getClient()

	req := CmdClientRequest{
		Name:        "randomname",
		ImageURL:    "pierone.stups.zalan.do/cat/cat-hello-aws:0.0.1",
		Replicas:    2,
		Ports:       []int{8080},
		Labels:      map[string]string{"name": "auto-test"},
		Env:         map[string]string{"name": "env-test"},
		CPULimit:    0,
		MemoryLimit: "2048",
		Force:       true,
	}

	clientino.CreateDeploy(&req)
	wait = true
}

func TestDeleteDeployment(t *testing.T) {
	clientino := getClient()
	clientino.InfoDeploy("randomname", true)
	if wait {
		fmt.Println("Please check successful deployment. Will now sleep for 5 minutes to make it end.")
		time.Sleep(20 * time.Second)
	}
	clientino.DeleteDeploy("randomname")
}

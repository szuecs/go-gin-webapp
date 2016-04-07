package client

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/resty.v0"

	konfig "github.com/szuecs/go-gin-webapp/conf/client"
	"golang.org/x/crypto/ssh/terminal"
)

//Client is the struct for accessing client functionalities
type Client struct {
	Config      *konfig.ClientConfig
	AccessToken string
}

var homeDirectories = []string{"HOME", "USERPROFILES"}
var tokenFilename = ".gin-oauth-token"

//RenewAccessToken is used to get a new OAuth2 access token
func (bc *Client) RenewAccessToken(username string) {
	if username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your username: ")
		username, _ = reader.ReadString('\n')
	}
	fmt.Print("Enter your password: ")
	bytePassword, err := terminal.ReadPassword(0)
	fmt.Println("")
	if err != nil {
		fmt.Printf("Cannot read password\n")
		os.Exit(1)
	}
	password := strings.TrimSpace(string(bytePassword))
	u, err := url.Parse(bc.Config.OauthURL)
	if err != nil {
		fmt.Printf("ERR: Could not parse given Auth URL: %s\n", bc.Config.OauthURL)
		os.Exit(1)
	}
	authURLStr := fmt.Sprintf("https://%s%s%s%s", u.Host, u.Path, u.RawQuery, u.Fragment)
	fmt.Printf("Getting token as %s\n", username)
	client := &http.Client{}
	req, err := http.NewRequest("GET", authURLStr, nil)
	req.SetBasicAuth(username, password)
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		fmt.Printf("ERR: Could not get Access Token, caused by: %s\n", err)
		os.Exit(1)
	}
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("ERR: Can not read response body, caused by: %s\n", err)
		os.Exit(1)
	}

	if len(respBody) > 0 && res.StatusCode == 200 {
		bc.AccessToken = string(respBody)
		fmt.Printf("SUCCESS. Your access token is stored in .chimp-token in your home directory.\n")
		//store token to file
		var homeDir string
		for _, home := range homeDirectories {
			if dir := os.Getenv(home); dir != "" {
				homeDir = dir
			}
		}
		tokenFileName := fmt.Sprintf("%s/%s", homeDir, tokenFilename)
		f, _ := os.Create(tokenFileName)
		_, _ = f.WriteString(strings.TrimSpace(bc.AccessToken)) //not important if doens't work, we'll try again next time
	} else {
		fmt.Printf("ERR: %d - %s\n", res.StatusCode, respBody)
	}
}

//GetAccessToken sets the access token inside the request
func (bc *Client) GetAccessToken(username string) {
	if bc.Config.Oauth2Enabled {
		//before trying to get the token I try to read the old one
		var homeDir string
		for _, home := range homeDirectories {
			if dir := os.Getenv(home); dir != "" {
				homeDir = dir
			}
		}
		tokenFileName := fmt.Sprintf("%s/%s", homeDir, tokenFilename)
		data, err := ioutil.ReadFile(tokenFileName)
		var oldToken string
		if err != nil {
			fmt.Println("ERR: Could not get an AccessToken which is required. Please login again.")
			os.Exit(1)
		} else {
			oldToken = strings.TrimSpace(string(data))
		}
		bc.AccessToken = oldToken
	}
}

// Get does HTTP GET to targetURL and print the result to STDOUT.
func (bc *Client) Get(targetURL *url.URL) {
	var resp *resty.Response
	var err error

	if bc.Config.Oauth2Enabled {
		resp, err = resty.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", bc.AccessToken)).
			Get(targetURL.String())
	} else {
		resp, err = resty.R().Get(targetURL.String())
	}
	if err != nil {
		fmt.Printf("ERR: Could not GET request, caused by: %s\n", err)
		os.Exit(1)
	}
	fmt.Print(resp)
}

// Package zalando contains Zalando specific definitions for
// authorization.
package zalando

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/zalando/gin-oauth2"
	"golang.org/x/oauth2"
)

// AccessTuples has to be set by the client to grant access.
var AccessTuples []AccessTuple

// AccessTuple is the type defined for use in AccessTuples.
type AccessTuple struct {
	Realm string `yaml:"Realm,omitempty"` // p.e. "employees", "services"
	Uid   string `yaml:"Uid,omitempty"`   // UnixName
	Cn    string `yaml:"Cn,omitempty"`    // RealName
}

// TeamInfo is defined like in TeamAPI json.
type TeamInfo struct {
	Id      string
	Id_name string
	Team_id string
	Type    string
	Name    string
	Mail    []string
}

// OAuth2Endpoint is similar to the definitions in golang.org/x/oauth2
var OAuth2Endpoint = oauth2.Endpoint{
	AuthURL:  "https://token.auth.zalando.com/access_token",
	TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
}

// TeamAPI is a custom API
var TeamAPI string = "https://teams.auth.zalando.com/api/teams"

// RequestTeamInfo is a function that returns team information for a
// given token.
func RequestTeamInfo(tc *ginoauth2.TokenContainer, uri string) ([]byte, error) {
	var uv = make(url.Values)
	uv.Set("member", tc.Scopes["uid"].(string))
	info_url := uri + "?" + uv.Encode()
	client := &http.Client{Transport: &ginoauth2.Transport}
	req, err := http.NewRequest("GET", info_url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tc.Token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// GroupCheck is an authorization function that checks, if the Token
// was issued for an employee of a specified team. The given
// TokenContainer must be valid. As side effect it sets "uid" and
// "team" in the gin.Context to the "official" team.
func GroupCheck(at []AccessTuple) func(tc *ginoauth2.TokenContainer, ctx *gin.Context) bool {
	ats := at
	return func(tc *ginoauth2.TokenContainer, ctx *gin.Context) bool {

		blob, err := RequestTeamInfo(tc, TeamAPI)
		if err != nil {
			glog.Error("failed to get team info, caused by: ", err)
			return false
		}
		var data []TeamInfo
		err = json.Unmarshal(blob, &data)
		if err != nil {
			glog.Errorf("JSON.Unmarshal failed, caused by: %s", err)
			return false
		}
		granted := false
		for _, teamInfo := range data {
			for idx := range ats {
				at := ats[idx]
				if teamInfo.Id == at.Uid {
					granted = true
					glog.Infof("Grant access to %s as team member of \"%s\"\n", tc.Scopes["uid"].(string), teamInfo.Id)
				}
				if teamInfo.Type == "official" {
					ctx.Set("uid", tc.Scopes["uid"].(string))
					ctx.Set("team", teamInfo.Id)
				}
			}
		}
		return granted
	}
}

// UidCheck is an authorization function that checks UID scope
// TokenContainer must be Valid. As side effect it sets "uid" and
// "cn" in the gin.Context to the authorized uid and cn (Realname).
func UidCheck(at []AccessTuple) func(tc *ginoauth2.TokenContainer, ctx *gin.Context) bool {
	ats := at
	return func(tc *ginoauth2.TokenContainer, ctx *gin.Context) bool {
		uid := tc.Scopes["uid"].(string)
		for idx := range ats {
			at := ats[idx]
			if tc.Realm == at.Realm && uid == at.Uid {
				ctx.Set("uid", uid)  //in this way I can set the authorized uid
				ctx.Set("cn", at.Cn) //in this way I can set the authorized Realname
				glog.Infof("Grant access to %s\n", uid)
				return true
			}
		}

		return false
	}
}

// NoAuthorization sets "team" and "uid" in the context without
// checking if the user/team is authorized.
func NoAuthorization(tc *ginoauth2.TokenContainer, ctx *gin.Context) bool {
	blob, err := RequestTeamInfo(tc, TeamAPI)
	var data []TeamInfo
	err = json.Unmarshal(blob, &data)
	if err != nil {
		glog.Errorf("JSON.Unmarshal failed, caused by: %s", err)
	}
	for _, teamInfo := range data {
		if teamInfo.Type == "official" {
			ctx.Set("uid", tc.Scopes["uid"].(string))
			ctx.Set("team", teamInfo.Id)
			return true
		}
	}
	return true
}

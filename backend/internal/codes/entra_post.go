package codes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/z3vxo/TokenJack/internal/db"
)

var ResourceToNameMap = map[string]string{
	"graph": "https://graph.microsoft.com/",
}

var httpClient = &http.Client{Timeout: 15 * time.Second}

type EntraResp struct {
	UserCode   string `json:"user_code"`
	DeviceCode string `json:"device_code"`
	VerifyUrl  string `json:"verification_uri"`
	Expires    int    `json:"expires_in"`
	Interval   int    `json:"interval"`
	Tenant     string
}

func (c *Codes) PostEntraDeviceCode(info db.CodeInfo) (EntraResp, error) {
	Url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode", info.TenantId)
	ClientId := "04b07795-8ddb-461a-bbee-02f9e1bf7b46" // hard coded value for now, just for testing

	var Scopes []string
	Scopes = append(Scopes, "openid")
	Scopes = append(Scopes, "profile")
	Scopes = append(Scopes, "offline_access") // remove later, dont need refresh for app
	for _, scope := range info.Scopes {
		Scopes = append(Scopes, fmt.Sprintf("%s%s", ResourceToNameMap[info.Resource], scope))
	}
	ScopeString := strings.Join(Scopes, " ")

	data := url.Values{}
	data.Set("client_id", ClientId)
	data.Set("scope", ScopeString)

	resp, err := httpClient.PostForm(Url, data)
	if err != nil {
		return EntraResp{}, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return EntraResp{}, fmt.Errorf("device code request failed: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var dcr EntraResp
	if err := json.Unmarshal(body, &dcr); err != nil {
		return EntraResp{}, err
	}
	dcr.Tenant = info.TenantId

	return dcr, nil
}

func (c *Codes) PollDeviceCode(dcr EntraResp) error {
	interbal := time.Duration(dcr.Interval) * time.Second
	tick := time.NewTicker(interbal)
	defer tick.Stop()

	deadline := time.Now().Add(time.Duration(dcr.Expires) * time.Second)

	for {
		select {
		case <-tick.C:
			if time.Now().After(deadline) {
				c.Remove(dcr)
				return fmt.Errorf("device code expired")
			}

			tokens, err := c.CheckToken(dcr)
			if err == nil {
				c.DB.InsertTokens(dcr.Tenant, tokens.AccessToken, tokens.RefreshToken)
				c.Remove(dcr)
				return nil
			}

		}

	}

}

type TokenResult struct {
	AccessToken  string
	RefreshToken string
}

func (c *Codes) CheckToken(dcr EntraResp) (TokenResult, error) {
	Url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", dcr.Tenant)
	ClientId := "04b07795-8ddb-461a-bbee-02f9e1bf7b46"

	resp, err := httpClient.PostForm(Url, url.Values{
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"client_id":   {ClientId},
		"device_code": {dcr.DeviceCode},
	})
	if err != nil {
		return TokenResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TokenResult{}, fmt.Errorf("authorization_pending")
	}
	var res map[string]any
	json.NewDecoder(resp.Body).Decode(&res)
	a, _ := res["access_token"].(string)
	r, _ := res["refresh_token"].(string)

	return TokenResult{
		AccessToken:  a,
		RefreshToken: r,
	}, nil
}

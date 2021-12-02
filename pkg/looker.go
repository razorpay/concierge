package pkg

import (
	"concierge/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

var client *LookerClient

type LookerClient struct {
	accessToken         string
	accessTokenExpireAt int64
	baseUrl             string
	clientId            string
	clientSecret        string
	datumUrl            string
	datumSecret         string
}

type LookerPatchUserRequest struct {
	IsDisabled bool `json:"is_disabled"`
}

type LookerPatchUserResponse struct {
	IsDisabled bool `json:"is_disabled"`
}

type LookerSearchUserRequest struct {
	Email string `json:"email"`
}

type LookerSearchUserResponse struct {
	Id         int  `json:"id"`
	IsDisabled bool `json:"is_disabled"`
}

func GetLookerClient() *LookerClient {
	if client == nil {
		client = &LookerClient{
			accessToken:         "",
			accessTokenExpireAt: 0,
			baseUrl:             config.LookerConfig.BaseUrl,
			clientId:            config.LookerConfig.ClientId,
			clientSecret:        config.LookerConfig.ClientSecret,
			datumUrl:            config.LookerConfig.DatumHostname,
			datumSecret:         config.LookerConfig.DatumAuthSecret,
		}
	}
	return client
}

func (c *LookerClient) PatchUser(userId int, req LookerPatchUserRequest) (*LookerPatchUserResponse, error) {
	path := fmt.Sprintf("api/3.1/users/%d", userId)
	method := http.MethodPatch
	body := struct {
		IsDisabled bool `json:"is_disabled"`
	}{req.IsDisabled}

	resp := LookerPatchUserResponse{}

	httpResponse, httpErr := c.executeRequest(path, method, body)

	if httpErr != nil {
		return nil, httpErr
	}

	if decodeErr := json.NewDecoder(httpResponse.Body).Decode(&resp); decodeErr != nil {
		return nil, decodeErr
	}
	return &resp, nil
}

func (c *LookerClient) SearchUser(req LookerSearchUserRequest) ([]LookerSearchUserResponse, error) {
	var resp []LookerSearchUserResponse

	path := fmt.Sprintf("api/3.1/users/search?email=%s", url.QueryEscape(req.Email))

	httpResponse, httpErr := c.executeRequest(path, http.MethodGet, nil)

	if httpErr != nil {
		return nil, httpErr
	}

	if decodeErr := json.NewDecoder(httpResponse.Body).Decode(&resp); decodeErr != nil {
		return nil, decodeErr
	}

	return resp, nil
}

func (c *LookerClient) UpdateLookerUserAttribute(email string) error {

	url := c.datumUrl + "/looker/qubole_api_token_update"
	body := struct {
		Email string `json:"email"`
	}{email}

	req := gorequest.New()
	req = req.Post(url).Send(body)

	req.Header.Set("X-Auth-Token", c.datumSecret)

	resp, _, errs := req.End()

	if resp.StatusCode != 200 {
		return errors.New("there was an error updating Looker user attribute. Please try Opening Looker once. " +
			"If dashboards are inaccessible, please contact `data platform` team")
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (c *LookerClient) CreateLookerUser(email string) error {

	url := c.datumUrl + "/looker/onboard-user"
	body := struct {
		Email string `json:"email"`
	}{email}

	req := gorequest.New()
	req = req.Post(url).Send(body)

	req.Header.Set("X-Auth-Token", c.datumSecret)

	resp, _, errs := req.End()

	if resp.StatusCode != 200 {
		return errors.New("there was an error while creating Looker user. Please contact `data platform` team")
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (c *LookerClient) isAccessTokenExpired() bool {
	return time.Now().Unix() >= c.accessTokenExpireAt
}

func (c *LookerClient) refreshAccessToken() error {
	requestBody := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{c.clientId, c.clientSecret}

	requestJson, _ := json.Marshal(&requestBody)

	requestString := string(requestJson)

	httpResponse, _, errs := gorequest.New().
		Post(c.baseUrl + "login").
		Type(gorequest.TypeForm).
		Send(requestString).
		End()

	if errs != nil {
		return errs[0]
	}

	if err := getErrorFromResponseIfApplicable(httpResponse); err != nil {
		return err
	}

	response := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}{}

	if decodeErr := json.NewDecoder(httpResponse.Body).Decode(&response); decodeErr != nil {
		return decodeErr
	}

	c.accessToken = response.AccessToken
	c.accessTokenExpireAt = time.Now().Unix() + response.ExpiresIn

	return nil
}

func (c *LookerClient) executeRequest(path string, method string, body interface{}) (*http.Response, error) {
	if c.isAccessTokenExpired() {
		if err := c.refreshAccessToken(); err != nil {
			return nil, err
		}
	}

	url := c.baseUrl + path

	req := gorequest.New()

	switch method {
	case http.MethodGet:
		req = req.Get(url)
	case http.MethodPatch:
		req = req.Patch(url).Send(body)
	case http.MethodPost:
		req = req.Post(url).Send(body)
	}

	req.Header.Set("Authorization", "Bearer token "+c.accessToken)

	resp, _, errs := req.End()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	if err := getErrorFromResponseIfApplicable(resp); err != nil {
		return nil, err
	}

	return resp, nil

}

func getErrorFromResponseIfApplicable(response gorequest.Response) error {
	if response.StatusCode >= 500 {
		return errors.New("there was an error contacting looker. please try after sometime. If error persists, " +
			"please contact looker admins")
	}

	if response.StatusCode >= 400 {
		return errors.New("failed to perform action on looker. please contact concierge admins if error persists")
	}

	return nil
}

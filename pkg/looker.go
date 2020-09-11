package pkg

import (
	"bytes"
	"concierge/config"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var client *LookerClient

type LookerClient struct {
	accessToken         string
	accessTokenExpireAt int64
	baseUrl             string
	clientId            string
	clientSecret        string
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

func (c *LookerClient) isAccessTokenExpired() bool {
	return time.Now().Unix() >= c.accessTokenExpireAt
}

func (c *LookerClient) setAccessToken() error {
	path := "login"
	method := http.MethodPost
	body := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_string"`
	}{c.clientId, c.clientSecret}

	httpResponse, httpErr := c.makeRequestAndResponse(path, method, body, nil)

	if httpErr != nil {
		return httpErr
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
		if err := c.setAccessToken(); err != nil {
			return nil, err
		}
	}

	headers := map[string]string{
		"Authoriziation": "Bearer token " + c.accessToken,
	}

	return c.makeRequestAndResponse(path, method, body, headers)
}

func (c *LookerClient) makeRequestAndResponse(path string, method string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := c.baseUrl + path

	switch method {
	case http.MethodPost:
		fallthrough
	case http.MethodPatch:
		var bodyJson []byte
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		return http.Post(url, "application/json", bytes.NewReader(bodyJson))
	case http.MethodGet:
		return http.Get(url)

	}

	return nil, nil
}

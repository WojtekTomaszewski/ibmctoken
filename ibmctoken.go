// Package ibmctoken gets IBM Cloud Oauth access token for provided API key.
//
// Usage
//
//   apiKey := os.Getenv("MY_API_KEY")
//   token := ibmctoken.NewToken(apiKey)
//   _ = token.RequestToken()
//   fmt.Println(token.AccessToken)
package ibmctoken

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// token endpoint address.
const (
	tokenUrl         = "https://iam.cloud.ibm.com/identity/token"
	tokenContentType = "application/x-www-form-urlencoded"
	tokenAccept      = "application/json"
	tokenGrantType   = "urn:ibm:params:oauth:grant-type:apikey"
)

// httpclient interface is used to mock http client.
type httpclient interface {
	Do(req *http.Request) (*http.Response, error)
}

// tokenResponse represents IBM Cloud Oauth access token response.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Expiration  int64  `json:"expiration"`
}

// tokenErrorResponse represents IBM Cloud Oauth access token error response.
type tokenErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
	ErrorDetails string `json:"errorDetails"`
}

// Token represents IBM Cloud Oauth access token.
type Token struct {
	ApiKey string
	Client httpclient
	tokenResponse
}

// NewToken creates new IBM Cloud Oauth access token.
func NewToken(apikey string) *Token {
	return &Token{
		ApiKey: apikey,
		Client: &http.Client{},
	}
}

// RequestToken fetches IBM Cloud Oauth access token for provided API key.
func (t *Token) RequestToken() error {
	return t.RequestTokenWithContext(context.Background())
}

// RequestTokenWithContext fetches IBM Cloud access token with custom context.
func (t *Token) RequestTokenWithContext(ctx context.Context) error {

	data := url.Values{}
	data.Set("grant_type", tokenGrantType)
	data.Set("apikey", t.ApiKey)

	dataEncoded := data.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenUrl, strings.NewReader(dataEncoded))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", tokenContentType)
	req.Header.Set("Accept", tokenAccept)

	res, err := t.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		errResponse := &tokenErrorResponse{}
		err := json.NewDecoder(res.Body).Decode(errResponse)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s %s", errResponse.ErrorMessage, errResponse.ErrorDetails)
	}

	if err := json.NewDecoder(res.Body).Decode(t); err != nil {
		return err
	}

	return nil
}

func (t *Token) Expired() bool {
	return time.Now().Unix() > t.Expiration
}

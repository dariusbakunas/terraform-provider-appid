package appid

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

type TokenResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken *string `json:"refresh_token,omitempty"` // when IAM api key used, this is never included
	Expiration   int64   `json:"expiration"`
	ExpiresIn    int64   `json:"expires_in"`
	Scope        string  `json:"scope"`
}

func getAccessToken(ctx context.Context, iamBaseURL string, apiKey string) (*TokenResponse, error) {
	log.Printf("[DEBUG] Getting IAM access token")

	baseURL, err := url.Parse(iamBaseURL)

	if err != nil {
		return nil, err
	}

	tokenURL, err := baseURL.Parse("/identity/token")

	if err != nil {
		return nil, err
	}

	c := &http.Client{
		Timeout: time.Minute * 2,
	}

	body := struct {
		ApiKey    string `url:"apikey"`
		GrantType string `url:"grant_type"`
	}{
		ApiKey:    apiKey,
		GrantType: "urn:ibm:params:oauth:grant-type:apikey",
	}

	token := &TokenResponse{}

	values, err := query.Values(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tokenURL.String(), strings.NewReader(values.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)

	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to get IAM access token")
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

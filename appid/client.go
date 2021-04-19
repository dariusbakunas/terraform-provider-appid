package appid

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	userAgent      = "go-terraform-provider-appid"
	defaultTimeout = time.Minute * 4
)

type service struct {
	client *Client
}

type Client struct {
	client    *http.Client // HTTP client used to communicate with the AppID API
	baseURL   *url.URL
	userAgent string

	Config *ConfigService
}

func NewClient(baseURL string, httpClient *http.Client) (*Client, error) {
	bURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(bURL.Path, "/") {
		bURL.Path += "/"
	}

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	c := &Client{client: httpClient, baseURL: bURL, userAgent: userAgent}

	baseService := service{c}

	c.Config = (*ConfigService)(&baseService)

	return c, nil
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	url, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	return req, nil
}

type AppIDError struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
}

func (r *AppIDError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c < 300 {
		return nil
	}

	errorResponse := &AppIDError{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		errorResponse.Message = string(data)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	// TODO: handle more status codes here

	return errorResponse
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)

	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)

	if err != nil {
		return resp, err
	}

	err = json.NewDecoder(resp.Body).Decode(v)

	if err == io.EOF {
		err = nil // ignore EOF errors caused by empty response body
	}

	return resp, err
}

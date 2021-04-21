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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2"

	"github.com/google/go-querystring/query"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"iam_api_key": {
				Type:        schema.TypeString,
				Description: "IBM Cloud IAM api key",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("IAM_API_KEY", nil),
			},
			"iam_access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("IAM_ACCESS_TOKEN", nil),
			},
			"appid_base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AppID API base URL, eg. https://us-south.appid.cloud.ibm.com",
				DefaultFunc: schema.EnvDefaultFunc("APPID_BASE_URL", nil),
			},
			"iam_base_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IBM IAM base URL, eg. https://iam.cloud.ibm.com",
				DefaultFunc: schema.EnvDefaultFunc("IAM_BASE_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"appid_token_config": resourceAppIDConfigTokens(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"appid_token_config": dataSourceAppIDConfigTokens(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type TokenResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken *string `json:"refresh_token,omitempty"` // when IAM api key used, this is never included
	Expiration   int64   `json:"expiration"`
	ExpiresIn    int64   `json:"expires_in"`
	Scope        string  `json:"scope"`
}

func getAccessToken(ctx context.Context, url string, apiKey string) (*TokenResponse, error) {
	log.Printf("[DEBUG] Getting IAM access token")

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

	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))

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

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var iamApiKey, iamAccesToken string

	if apiKey, ok := d.GetOk("iam_api_key"); ok {
		iamApiKey = apiKey.(string)
	}

	if accessToken, ok := d.GetOk("iam_access_token"); ok {
		iamAccesToken = accessToken.(string)
	}

	appIDBaseURL := d.Get("appid_base_url").(string)
	iamBaseURL, err := url.Parse(d.Get("iam_base_url").(string))

	if err != nil {
		return nil, diag.FromErr(err)
	}

	tokenURL, err := iamBaseURL.Parse("/identity/token")

	if err != nil {
		return nil, diag.FromErr(err)
	}

	if iamAccesToken == "" {
		if iamApiKey == "" {
			return nil, diag.Errorf("iam_api_key or iam_access_token must be specified")
		}

		token, err := getAccessToken(ctx, tokenURL.String(), iamApiKey)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		iamAccesToken = token.AccessToken
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: iamAccesToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	c, err := NewClient(appIDBaseURL, tc)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}

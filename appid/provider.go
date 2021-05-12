package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.ibm.com/dbakuna/terraform-provider-appid/api"
	"golang.org/x/oauth2"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"iam_api_key": {
				Type:        schema.TypeString,
				Description: "The IBM Cloud IAM api key used to retrieve IAM access token if `iam_access_token` is not specified",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("IAM_API_KEY", nil),
			},
			"iam_access_token": {
				Type:        schema.TypeString,
				Description: "The IBM Cloud Identity and Access Management token used to access AppID APIs",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("IAM_ACCESS_TOKEN", nil),
			},
			"appid_base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AppID API base URL (for example 'https://us-south.appid.cloud.ibm.com')",
				DefaultFunc: schema.EnvDefaultFunc("APPID_BASE_URL", nil),
			},
			"iam_base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IBM IAM base URL",
				Default:     "https://iam.cloud.ibm.com",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IBM cloud Region (for example 'us-south').",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_REGION", "IBMCLOUD_REGION"}, "us-south"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"appid_application":              resourceAppIDApplication(),
			"appid_cloud_directory_template": resourceAppIDCloudDirectoryTemplate(),
			"appid_idp_custom":               resourceAppIDIDPCustom(),
			"appid_idp_saml":                 resourceAppIDIDPSaml(),
			"appid_role":                     resourceAppIDRole(),
			"appid_token_config":             resourceAppIDTokenConfig(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"appid_application":              dataSourceAppIDApplication(),
			"appid_applications":             dataSourceAppIDApplications(),
			"appid_application_ids":          dataSourceAppIDApplicationIDs(),
			"appid_cloud_directory_template": dataSourceAppIDCloudDirectoryTemplate(),
			"appid_idp_cloud_directory":      dataSourceAppIDIDPCloudDirectory(),
			"appid_idp_custom":               dataSourceAppIDIDPCustom(),
			"appid_idp_saml":                 dataSourceAppIDIDPSAML(),
			"appid_role":                     dataSourceAppIDRole(),
			"appid_token_config":             dataSourceAppIDTokenConfig(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var iamApiKey, iamAccesToken string

	clientOptions := &api.Options{}

	region := d.Get("region").(string)
	baseURL := d.Get("appid_base_url").(string)

	if region == "" && baseURL == "" {
		return nil, diag.Errorf("region or baseURL must be specified")
	}

	if region != "" && baseURL != "" {
		log.Printf("[WARN] both region and baseURL were specified, baseURL will take precendence")
	}

	if apiKey, ok := d.GetOk("iam_api_key"); ok {
		iamApiKey = apiKey.(string)
	}

	if accessToken, ok := d.GetOk("iam_access_token"); ok {
		iamAccesToken = accessToken.(string)
	}

	if region != "" {
		clientOptions.BaseURL = fmt.Sprintf("https://%s.appid.cloud.ibm.com", region)
	}

	if baseURL != "" {
		clientOptions.BaseURL = baseURL
	}

	if iamAccesToken == "" {
		if iamApiKey == "" {
			return nil, diag.Errorf("iam_api_key or iam_access_token must be specified")
		}

		token, err := getAccessToken(ctx, d.Get("iam_base_url").(string), iamApiKey)

		if err != nil {
			return nil, diag.FromErr(err)
		}

		iamAccesToken = token.AccessToken
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: iamAccesToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	c, err := api.NewClient(clientOptions, tc)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}

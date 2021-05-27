package appid

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path"

	appid "github.com/IBM/appid-go-sdk/appidmanagementv4"
	"github.com/IBM/go-sdk-core/core"

	//v5core "github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_API_KEY", "IBMCLOUD_API_KEY", "IAM_API_KEY"}, nil),
			},
			"iam_access_token": {
				Type:        schema.TypeString,
				Description: "The IBM Cloud Identity and Access Management token used to access AppID APIs",
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"IC_IAM_TOKEN", "IBMCLOUD_IAM_TOKEN", "IAM_ACCESS_TOKEN"}, nil),
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
			"api_max_retry": {
				Description: "Maximum number of retries for AppID api requests, set to 0 to disable",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"appid_action_url":               resourceAppIDActionURL(),
			"appid_application":              resourceAppIDApplication(),
			"appid_audit_status":             resourceAppIDAuditStatus(),
			"appid_cloud_directory_template": resourceAppIDCloudDirectoryTemplate(),
			"appid_idp_cloud_directory":      resourceAppIDIDPCloudDirectory(),
			"appid_idp_custom":               resourceAppIDIDPCustom(),
			"appid_idp_facebook":             resourceAppIDIDPFacebook(),
			"appid_idp_google":               resourceAppIDIDPGoogle(),
			"appid_idp_saml":                 resourceAppIDIDPSaml(),
			"appid_media":                    resourceAppIDMedia(),
			"appid_password_regex":           resourceAppIDPasswordRegex(),
			"appid_redirect_urls":            resourceAppIDRedirectURLs(),
			"appid_role":                     resourceAppIDRole(),
			"appid_theme_color":              resourceAppIDThemeColor(),
			"appid_token_config":             resourceAppIDTokenConfig(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"appid_action_url":               dataSourceAppIDActionURL(),
			"appid_application":              dataSourceAppIDApplication(),
			"appid_applications":             dataSourceAppIDApplications(),
			"appid_application_ids":          dataSourceAppIDApplicationIDs(),
			"appid_audit_status":             dataSourceAppIDAuditStatus(),
			"appid_cloud_directory_template": dataSourceAppIDCloudDirectoryTemplate(),
			"appid_password_regex":           dataSourceAppIDPasswordRegex(),
			"appid_idp_cloud_directory":      dataSourceAppIDIDPCloudDirectory(),
			"appid_idp_custom":               dataSourceAppIDIDPCustom(),
			"appid_idp_facebook":             dataSourceAppIDIDPFacebook(),
			"appid_idp_google":               dataSourceAppIDIDPGoogle(),
			"appid_idp_saml":                 dataSourceAppIDIDPSAML(),
			"appid_media":                    dataSourceAppIDMedia(),
			"appid_redirect_urls":            dataSourceAppIDRedirectURLs(),
			"appid_role":                     dataSourceAppIDRole(),
			"appid_roles":                    dataSourceAppIDRoles(),
			"appid_theme_color":              dataSourceAppIDThemeColor(),
			"appid_token_config":             dataSourceAppIDTokenConfig(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var iamApiKey, iamAccessToken string

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
		iamAccessToken = accessToken.(string)
	}

	options := &appid.AppIDManagementV4Options{}

	if region != "" {
		options.URL = fmt.Sprintf("https://%s.appid.cloud.ibm.com", region)
	}

	if baseURL != "" {
		options.URL = baseURL
	}

	if iamAccessToken == "" {
		if iamApiKey == "" {
			return nil, diag.Errorf("iam_api_key or iam_access_token must be specified")
		}

		iamBaseURL := d.Get("iam_base_url").(string)

		u, err := url.Parse(iamBaseURL)

		if err != nil {
			return nil, diag.Errorf("failed parsing iam_base_url")
		}

		u.Path = path.Join(u.Path, "/identity/token")

		options.Authenticator = &core.IamAuthenticator{
			ApiKey: iamApiKey,
			URL:    u.String(),
		}
	} else {
		options.Authenticator = &core.BearerTokenAuthenticator{
			BearerToken: iamAccessToken,
		}
	}

	//v5core.GetLogger().SetLogLevel(v5core.LevelDebug)
	client, err := appid.NewAppIDManagementV4(options)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.EnableRetries(d.Get("api_max_retry").(int), 0) // 0 delay - using client default
	return client, diags
}

package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xilution/xilution-client-go"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_ORGANIZATION_ID", nil),
			},
			"grant_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_GRANT_TYPE", "client_credentials"),
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_SCOPE", "read write"),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_CLIENT_SECRET", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_PASSWORD", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"xilution_organization":                  dataSourceOrganization(),
			"xilution_client":                        dataSourceClient(),
			"xilution_user":                          dataSourceUser(),
			"xilution_git_account":                   dataSourceGitAccount(),
			"xilution_git_repo":                      dataSourceGitRepo(),
			"xilution_git_repo_event":                dataSourceGitRepoEvent(),
			"xilution_cloud_provider":                dataSourceCloudProvider(),
			"xilution_vpc_pipeline":                  dataSourceVpcPipeline(),
			"xilution_vpc_pipeline_event":            dataSourceVpcPipelineEvent(),
			"xilution_k8s_pipeline":                  dataSourceK8sPipeline(),
			"xilution_k8s_pipeline_event":            dataSourceK8sPipelineEvent(),
			"xilution_word_press_pipeline":           dataSourceWordPressPipeline(),
			"xilution_word_press_pipeline_event":     dataSourceWordPressPipelineEvent(),
			"xilution_static_content_pipeline":       dataSourceStaticContentPipeline(),
			"xilution_static_content_pipeline_event": dataSourceStaticContentPipelineEvent(),
			"xilution_api_pipeline":                  dataSourceApiPipeline(),
			"xilution_api_pipeline_event":            dataSourceApiPipelineEvent(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"xilution_git_account":                   resourceGitAccount(),
			"xilution_git_repo":                      resourceGitRepo(),
			"xilution_git_repo_event":                resourceGitRepoEvent(),
			"xilution_cloud_provider":                resourceCloudProvider(),
			"xilution_vpc_pipeline":                  resourceVpcPipeline(),
			"xilution_vpc_pipeline_event":            resourceVpcPipelineEvent(),
			"xilution_k8s_pipeline":                  resourceK8sPipeline(),
			"xilution_k8s_pipeline_event":            resourceK8sPipelineEvent(),
			"xilution_word_press_pipeline":           resourceWordPressPipeline(),
			"xilution_word_press_pipeline_event":     resourceWordPressPipelineEvent(),
			"xilution_static_content_pipeline":       resourceStaticContentPipeline(),
			"xilution_static_content_pipeline_event": resourceStaticContentPipelineEvent(),
			"xilution_api_pipeline":                  resourceApiPipeline(),
			"xilution_api_pipeline_event":            resourceApiPipelineEvent(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Println("[INFO] Configuring Xilution Provider")

	organizationId := d.Get("organization_id").(string)
	grantType := d.Get("grant_type").(string)
	scope := d.Get("scope").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	xc, err := xilution.NewXilutionClient(&organizationId, &grantType, &scope, &clientId, &clientSecret, &username, &password)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Xilution client",
			Detail:   "Unable to auth user for authenticated Xilution client",
		})

		return nil, diags
	}

	return xc, diags
}

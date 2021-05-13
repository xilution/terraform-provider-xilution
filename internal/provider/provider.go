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
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_CLIENT_ID", nil),
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("XILUTION_ORGANIZATION_ID", nil),
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
			"xilution_organization":            dataSourceOrganization(),
			"xilution_client":                  dataSourceClient(),
			"xilution_user":                    dataSourceUser(),
			"xilution_git_account":             dataSourceGitAccount(),
			"xilution_git_repo":                dataSourceGitRepo(),
			"xilution_git_repo_event":          dataSourceGitRepoEvent(),
			"xilution_cloud_provider":          dataSourceCloudProvider(),
			"xilution_vpc_pipeline":            dataSourceVpcPipeline(),
			"xilution_k8s_pipeline":            dataSourceK8sPipeline(),
			"xilution_word_press_pipeline":     dataSourceWordPressPipeline(),
			"xilution_static_content_pipeline": dataSourceStaticContentPipeline(),
			"xilution_api_pipeline":            dataSourceApiPipeline(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"xilution_git_account":             resourceGitAccount(),
			"xilution_git_repo":                resourceGitRepo(),
			"xilution_git_repo_event":          resourceGitRepoEvent(),
			"xilution_cloud_provider":          resourceCloudProvider(),
			"xilution_vpc_pipeline":            resourceVpcPipeline(),
			"xilution_k8s_pipeline":            resourceK8sPipeline(),
			"xilution_word_press_pipeline":     resourceWordPressPipeline(),
			"xilution_static_content_pipeline": resourceStaticContentPipeline(),
			"xilution_api_pipeline":            resourceApiPipeline(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Println("[INFO] Configuring Xilution Provider")

	clientId := d.Get("client_id").(string)
	organizationId := d.Get("organization_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if clientId != "" && organizationId != "" && username != "" && password != "" {
		xc, err := xilution.NewXilutionClient(&clientId, &organizationId, &username, &password)
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

	xc, err := xilution.NewXilutionClient(nil, nil, nil, nil)
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

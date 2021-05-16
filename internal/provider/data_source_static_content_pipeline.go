package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceStaticContentPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStaticContentPipelineRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipeline_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_repo_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stages": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
				Optional: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owning_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStaticContentPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	staticContentPipelineId := d.Get("id").(string)

	staticContentPipeline, err := c.GetStaticContentPipeline(&organizationId, &staticContentPipelineId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", staticContentPipeline.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", staticContentPipeline.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_type", staticContentPipeline.PipelineType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cloud_provider_id", staticContentPipeline.CloudProviderId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_repo_id", staticContentPipeline.GitRepoId); err != nil {
		return diag.FromErr(err)
	}
	
	stages := make([]interface{}, len(staticContentPipeline.Stages))
	for i, stage := range staticContentPipeline.Stages {
		newStage := make(map[string]interface{})

		newStage["name"] = stage.Name
		stages[i] = newStage
	}
	if err := d.Set("stages", stages); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", staticContentPipeline.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", staticContentPipeline.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", staticContentPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", staticContentPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(staticContentPipeline.ID)

	return diags
}

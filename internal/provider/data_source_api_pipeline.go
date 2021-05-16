package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceApiPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiPipelineRead,
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
			"vpc_pipeline_id": {
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

func dataSourceApiPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	apiPipelineId := d.Get("id").(string)

	apiPipeline, err := c.GetApiPipeline(&organizationId, &apiPipelineId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", apiPipeline.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", apiPipeline.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_type", apiPipeline.PipelineType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vpc_pipeline_id", apiPipeline.VpcPipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_repo_id", apiPipeline.GitRepoId); err != nil {
		return diag.FromErr(err)
	}

	stages := make([]interface{}, len(apiPipeline.Stages))
	for i, stage := range apiPipeline.Stages {
		newStage := make(map[string]interface{})

		newStage["name"] = stage.Name
		stages[i] = newStage
	}
	if err := d.Set("stages", stages); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", apiPipeline.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", apiPipeline.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", apiPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", apiPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(apiPipeline.ID)

	return diags
}

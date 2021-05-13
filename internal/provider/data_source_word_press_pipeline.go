package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceWordPressPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWordPressPipelineRead,
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
			"k8s_pipeline_id": {
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
							Required: true,
						},
					},
				},
				Required: true,
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

func dataSourceWordPressPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	wordPressPipelineId := d.Get("id").(string)

	wordPressPipeline, err := c.GetWordPressPipeline(&organizationId, &wordPressPipelineId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", wordPressPipeline.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", wordPressPipeline.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_type", wordPressPipeline.PipelineType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("k8s_pipeline_id", wordPressPipeline.K8sPipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_repo_id", wordPressPipeline.GitRepoId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("stages", wordPressPipeline.Stages); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", wordPressPipeline.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", wordPressPipeline.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", wordPressPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", wordPressPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(wordPressPipeline.ID)

	return diags
}

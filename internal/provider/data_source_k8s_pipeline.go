package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceK8sPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceK8sPipelineRead,
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

func dataSourceK8sPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	k8sPipelineId := d.Get("id").(string)

	k8sPipeline, err := c.GetK8sPipeline(&organizationId, &k8sPipelineId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", k8sPipeline.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", k8sPipeline.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_type", k8sPipeline.PipelineType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vpc_pipeline_id", k8sPipeline.VpcPipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", k8sPipeline.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", k8sPipeline.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", k8sPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", k8sPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(k8sPipeline.ID)

	return diags
}

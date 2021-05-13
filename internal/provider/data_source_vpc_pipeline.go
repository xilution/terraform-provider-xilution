package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceVpcPipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcPipelineRead,
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

func dataSourceVpcPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	vpcPipelineId := d.Get("id").(string)

	vpcPipeline, err := c.GetVpcPipeline(&organizationId, &vpcPipelineId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", vpcPipeline.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", vpcPipeline.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_type", vpcPipeline.PipelineType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cloud_provider_id", vpcPipeline.CloudProviderId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", vpcPipeline.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", vpcPipeline.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", vpcPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", vpcPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vpcPipeline.ID)

	return diags
}

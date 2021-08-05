package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceApiPipelineEvent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiPipelineEventRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_type": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceApiPipelineEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	pipelineEventId := d.Get("id").(string)

	pipelineEvent, err := c.GetApiPipelineEvent(&organizationId, &pipelineEventId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", pipelineEvent.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", pipelineEvent.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_id", pipelineEvent.PipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_type", pipelineEvent.EventType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", pipelineEvent.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", pipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", pipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pipelineEvent.ID)

	return diags
}

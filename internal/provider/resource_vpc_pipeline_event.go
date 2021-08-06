package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceVpcPipelineEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcPipelineEventCreate,
		ReadContext:   resourceVpcPipelineEventRead,
		UpdateContext: resourceVpcPipelineEventUpdate,
		DeleteContext: resourceVpcPipelineEventDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owning_user_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceVpcPipelineEventCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	owningUserId := d.Get("owning_user_id").(string)
	eventType := d.Get("event_type").(string)

	location, err := c.CreateVpcPipelineEvent(&organizationId, &xc.PipelineEvent{
		Type:           "pipeline-event",
		PipelineId:     pipelineId,
		OrganizationId: organizationId,
		OwningUserId:   owningUserId,
		EventType:      eventType,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	time.Sleep(5 * time.Second)

	id := getIdFromLocationUrl(location)

	d.SetId(*id)

	getPipelineStatusFunc := func() (*xc.PipelineStatus, error) {
		pipeline, err := c.GetVpcPipeline(&organizationId, &pipelineId)
		if err != nil {
			return nil, err
		}
		return pipeline.Status, nil
	}

	err = waitForPipelineEventToComplete(eventType, 15*time.Minute, 5*time.Second, getPipelineStatusFunc)
	if err != nil {
		return diag.FromErr(err)
	}

	vpcPipelineEvent, err := c.GetVpcPipelineEvent(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", vpcPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", vpcPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceVpcPipelineEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	vpcPipelineEvent, err := c.GetVpcPipelineEvent(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", vpcPipelineEvent.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_id", vpcPipelineEvent.PipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", vpcPipelineEvent.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_type", vpcPipelineEvent.EventType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", vpcPipelineEvent.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", vpcPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", vpcPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceVpcPipelineEventUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceVpcPipelineEventRead(ctx, d, m)
}

func resourceVpcPipelineEventDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

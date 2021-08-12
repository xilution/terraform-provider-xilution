package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceApiPipelineEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiPipelineEventCreate,
		ReadContext:   resourceApiPipelineEventRead,
		UpdateContext: resourceApiPipelineEventUpdate,
		DeleteContext: resourceApiPipelineEventDelete,
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

func resourceApiPipelineEventCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	owningUserId := d.Get("owning_user_id").(string)
	eventType := d.Get("event_type").(string)

	location, err := c.CreateApiPipelineEvent(&organizationId, &xc.PipelineEvent{
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
		pipeline, err := c.GetApiPipeline(&organizationId, &pipelineId)
		if err != nil {
			return nil, err
		}
		return pipeline.Status, nil
	}

	err = waitForPipelineEventToComplete(eventType, 30*time.Minute, 5*time.Second, getPipelineStatusFunc)
	if err != nil {
		return diag.FromErr(err)
	}

	apiPipelineEvent, err := c.GetApiPipelineEvent(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", apiPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", apiPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceApiPipelineEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	apiPipelineEvent, err := c.GetApiPipelineEvent(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", apiPipelineEvent.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_id", apiPipelineEvent.PipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", apiPipelineEvent.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_type", apiPipelineEvent.EventType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", apiPipelineEvent.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", apiPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", apiPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceApiPipelineEventUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceApiPipelineEventRead(ctx, d, m)
}

func resourceApiPipelineEventDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

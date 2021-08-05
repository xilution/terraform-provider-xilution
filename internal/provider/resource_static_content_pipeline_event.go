package provider

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceStaticContentPipelineEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStaticContentPipelineEventCreate,
		ReadContext:   resourceStaticContentPipelineEventRead,
		UpdateContext: resourceStaticContentPipelineEventUpdate,
		DeleteContext: resourceStaticContentPipelineEventDelete,
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

func resourceStaticContentPipelineEventCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	eventType := d.Get("event_type").(string)

	location, err := c.CreateStaticContentPipelineEvent(&organizationId, &xc.PipelineEvent{
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

	timeoutInMinutes := 15.0
	done := false
	start := time.Now()
	for !done {
		pipeline, err := c.GetStaticContentPipeline(&organizationId, &pipelineId)
		if err != nil {
			return diag.FromErr(err)
		}
		status := pipeline.Status
		if pipeline.Status != nil {
			infrastructureStatus := status.InfrastructureStatus
			if infrastructureStatus == "CREATE_COMPLETE" {
				continuousIntegrationStatus := status.ContinuousIntegrationStatus
				if continuousIntegrationStatus != nil {
					latestUpExecutionStatus := continuousIntegrationStatus.LatestUpExecutionStatus
					if latestUpExecutionStatus == "SUCCEEDED" {
						done = true
					} else if latestUpExecutionStatus == "FAILED" {
						return diag.FromErr(errors.New("create pipeline event failed. pipeline up status is failed"))
					}
				}
			} else if infrastructureStatus == "CREATE_FAILED" {
				return diag.FromErr(errors.New("create pipeline event failed. pipeline infrastructure status is failed"))
			}
		} else {
			if time.Since(start).Minutes() > timeoutInMinutes {
				return diag.FromErr(err)
			}
			time.Sleep(5 * time.Second)
		}
	}

	staticcontentPipelineEvent, err := c.GetStaticContentPipelineEvent(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", staticcontentPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", staticcontentPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStaticContentPipelineEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	staticcontentPipelineEvent, err := c.GetStaticContentPipelineEvent(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", staticcontentPipelineEvent.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_id", staticcontentPipelineEvent.PipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", staticcontentPipelineEvent.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_type", staticcontentPipelineEvent.EventType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", staticcontentPipelineEvent.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", staticcontentPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", staticcontentPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStaticContentPipelineEventUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceStaticContentPipelineEventRead(ctx, d, m)
}

func resourceStaticContentPipelineEventDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

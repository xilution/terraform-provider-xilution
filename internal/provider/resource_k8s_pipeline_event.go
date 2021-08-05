package provider

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceK8sPipelineEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceK8sPipelineEventCreate,
		ReadContext:   resourceK8sPipelineEventRead,
		UpdateContext: resourceK8sPipelineEventUpdate,
		DeleteContext: resourceK8sPipelineEventDelete,
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

func resourceK8sPipelineEventCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	eventType := d.Get("event_type").(string)

	location, err := c.CreateK8sPipelineEvent(&organizationId, &xc.PipelineEvent{
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

	timeoutInMinutes := 45.0
	done := false
	start := time.Now()
	for !done {
		pipeline, err := c.GetK8sPipeline(&organizationId, &pipelineId)
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

	k8sPipelineEvent, err := c.GetK8sPipelineEvent(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", k8sPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", k8sPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceK8sPipelineEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	k8sPipelineEvent, err := c.GetK8sPipelineEvent(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", k8sPipelineEvent.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_id", k8sPipelineEvent.PipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", k8sPipelineEvent.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_type", k8sPipelineEvent.EventType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", k8sPipelineEvent.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", k8sPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", k8sPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceK8sPipelineEventUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceK8sPipelineEventRead(ctx, d, m)
}

func resourceK8sPipelineEventDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

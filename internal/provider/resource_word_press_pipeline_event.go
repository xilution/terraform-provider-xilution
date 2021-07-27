package provider

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceWordPressPipelineEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWordPressPipelineEventCreate,
		ReadContext:   resourceWordPressPipelineEventRead,
		UpdateContext: resourceWordPressPipelineEventUpdate,
		DeleteContext: resourceWordPressPipelineEventDelete,
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
			"parameters": {
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

func resourceWordPressPipelineEventCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	owningUserId := d.Get("owning_user_id").(string)
	parametersStr := d.Get("parameters").(string)
	var parametersData map[string]interface{}
	err := json.Unmarshal([]byte(parametersStr), &parametersData)
	if err != nil {
		return diag.FromErr(err)
	}

	eventType := d.Get("event_type").(string)

	location, err := c.CreateWordPressPipelineEvent(&organizationId, &xc.PipelineEvent{
		Type:           "pipeline-event",
		PipelineId:     pipelineId,
		OrganizationId: organizationId,
		OwningUserId:   owningUserId,
		Parameters:     parametersData,
		EventType:      eventType,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	id := getIdFromLocationUrl(location)

	d.SetId(*id)

	timeoutInMinutes := 15.0
	done := false
	start := time.Now()
	for !done {
		pipeline, err := c.GetWordPressPipeline(&organizationId, &pipelineId)
		if err != nil {
			return diag.FromErr(err)
		}
		status := pipeline.Status.ContinuousIntegrationStatus.LatestUpExecutionStatus
		log.Println("[DEBUG] WordPress Pipeline Status is ", status)
		if status == "SUCCEEDED" {
			done = true
		} else {
			if time.Since(start).Minutes() > timeoutInMinutes {
				return diag.FromErr(err)
			}
			time.Sleep(5 * time.Second)
		}
	}

	wordpressPipelineEvent, err := c.GetWordPressPipelineEvent(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", wordpressPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", wordpressPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWordPressPipelineEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	wordpressPipelineEvent, err := c.GetWordPressPipelineEvent(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", wordpressPipelineEvent.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_id", wordpressPipelineEvent.PipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", wordpressPipelineEvent.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	jsonStr, err := json.Marshal(wordpressPipelineEvent.Parameters)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] jsonStr: " + string(jsonStr))
	if err := d.Set("parameters", string(jsonStr)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_type", wordpressPipelineEvent.EventType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", wordpressPipelineEvent.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", wordpressPipelineEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", wordpressPipelineEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWordPressPipelineEventUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceWordPressPipelineEventRead(ctx, d, m)
}

func resourceWordPressPipelineEventDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

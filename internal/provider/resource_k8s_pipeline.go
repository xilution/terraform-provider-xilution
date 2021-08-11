package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceK8sPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceK8sPipelineCreate,
		ReadContext:   resourceK8sPipelineRead,
		UpdateContext: resourceK8sPipelineUpdate,
		DeleteContext: resourceK8sPipelineDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_pipeline_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
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

func resourceK8sPipelineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	vpcPipelineId := d.Get("vpc_pipeline_id").(string)
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	location, err := c.CreateK8sPipeline(&organizationId, &xc.K8sPipeline{
		Type:           "pipeline",
		Name:           name,
		PipelineType:   pipelineType,
		VpcPipelineId:  vpcPipelineId,
		OrganizationId: organizationId,
		OwningUserId:   owningUserId,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	id := getIdFromLocationUrl(location)

	d.SetId(*id)

	k8sPipeline, err := c.GetK8sPipeline(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", k8sPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", k8sPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceK8sPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	k8sPipeline, err := c.GetK8sPipeline(&organizationId, &id)
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

	return diags
}

func resourceK8sPipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	vpcPipelineId := d.Get("vpc_pipeline_id").(string)
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	if d.HasChange("name") {
		err := c.UpdateK8sPipeline(&organizationId, &xc.K8sPipeline{
			Type:           "pipeline",
			ID:             id,
			Name:           name,
			PipelineType:   pipelineType,
			VpcPipelineId:  vpcPipelineId,
			OrganizationId: organizationId,
			OwningUserId:   owningUserId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceK8sPipelineRead(ctx, d, m)
}

func resourceK8sPipelineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)
	id := d.Id()

	getPipelineStatusFunc := func() (*xc.PipelineStatus, error) {
		pipeline, err := GetK8sPipeline(&organizationId, &id)
		if err != nil {
			return nil, err
		}
		return pipeline.Status, nil
	}

	status, err := getPipelineStatusFunc()
	if err != nil {
		return diag.FromErr(err)
	}

	if status.InfrastructureStatus != NOT_FOUND {
		_, err = c.CreateK8sPipelineEvent(&organizationId, &xc.PipelineEvent{
			Type:           "pipeline-event",
			PipelineId:     id,
			OrganizationId: organizationId,
			OwningUserId:   owningUserId,
			EventType:      "DEPROVISION",
		})
		if err != nil {
			return diag.FromErr(err)
		}
		time.Sleep(5 * time.Second)

		err = waitForPipelineInfrastructureNotFound(45*time.Minute, 5*time.Second, getPipelineStatusFunc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = c.DeleteK8sPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

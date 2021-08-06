package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceVpcPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcPipelineCreate,
		ReadContext:   resourceVpcPipelineRead,
		UpdateContext: resourceVpcPipelineUpdate,
		DeleteContext: resourceVpcPipelineDelete,
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
			"cloud_provider_id": {
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

func resourceVpcPipelineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	cloudProviderId := d.Get("cloud_provider_id").(string)
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	location, err := c.CreateVpcPipeline(&organizationId, &xc.VpcPipeline{
		Type:            "pipeline",
		Name:            name,
		PipelineType:    pipelineType,
		CloudProviderId: cloudProviderId,
		OrganizationId:  organizationId,
		OwningUserId:    owningUserId,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	id := getIdFromLocationUrl(location)

	d.SetId(*id)

	vpcPipeline, err := c.GetVpcPipeline(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", vpcPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", vpcPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceVpcPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	vpcPipeline, err := c.GetVpcPipeline(&organizationId, &id)
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

	return diags
}

func resourceVpcPipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	cloudProviderId := d.Get("cloud_provider_id").(string)
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	if d.HasChange("name") {
		err := c.UpdateVpcPipeline(&organizationId, &xc.VpcPipeline{
			Type:            "pipeline",
			ID:              id,
			Name:            name,
			PipelineType:    pipelineType,
			CloudProviderId: cloudProviderId,
			OrganizationId:  organizationId,
			OwningUserId:    owningUserId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVpcPipelineRead(ctx, d, m)
}

func resourceVpcPipelineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)
	id := d.Id()

	getPipelineStatusFunc := func() (*xc.PipelineStatus, error) {
		pipeline, err := c.GetVpcPipeline(&organizationId, &id)
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
		_, err = c.CreateVpcPipelineEvent(&organizationId, &xc.PipelineEvent{
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

		err = waitForPipelineInfrastructureNotFound(15*time.Minute, 5*time.Second, getPipelineStatusFunc)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = c.DeleteVpcPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

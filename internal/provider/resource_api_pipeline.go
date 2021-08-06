package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceApiPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiPipelineCreate,
		ReadContext:   resourceApiPipelineRead,
		UpdateContext: resourceApiPipelineUpdate,
		DeleteContext: resourceApiPipelineDelete,
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
			"git_repo_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stages": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func resourceApiPipelineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	vpcPipelineId := d.Get("vpc_pipeline_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	branch := d.Get("branch").(string)
	stages := d.Get("stages").([]interface{})
	mappedStages := []xc.ApiStage{}
	for _, stage := range stages {
		newStage := xc.ApiStage{
			Name: stage.(map[string]interface{})["name"].(string),
		}
		mappedStages = append(mappedStages, newStage)
	}
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	location, err := c.CreateApiPipeline(&organizationId, &xc.ApiPipeline{
		Type:           "pipeline",
		Name:           name,
		PipelineType:   pipelineType,
		VpcPipelineId:  vpcPipelineId,
		GitRepoId:      gitRepoId,
		Branch:         branch,
		Stages:         mappedStages,
		OrganizationId: organizationId,
		OwningUserId:   owningUserId,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	id := getIdFromLocationUrl(location)

	d.SetId(*id)

	apiPipeline, err := c.GetApiPipeline(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", apiPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", apiPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceApiPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	apiPipeline, err := c.GetApiPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", apiPipeline.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", apiPipeline.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_type", apiPipeline.PipelineType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("vpc_pipeline_id", apiPipeline.VpcPipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_repo_id", apiPipeline.GitRepoId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("branch", apiPipeline.Branch); err != nil {
		return diag.FromErr(err)
	}

	stages := make([]interface{}, len(apiPipeline.Stages))
	for i, stage := range apiPipeline.Stages {
		newStage := make(map[string]interface{})

		newStage["name"] = stage.Name
		stages[i] = newStage
	}
	if err := d.Set("stages", stages); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", apiPipeline.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", apiPipeline.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", apiPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", apiPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceApiPipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	vpcPipelineId := d.Get("vpc_pipeline_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	branch := d.Get("branch").(string)
	stages := d.Get("stages").([]interface{})
	mappedStages := []xc.ApiStage{}
	for _, stage := range stages {
		newStage := xc.ApiStage{
			Name: stage.(map[string]interface{})["name"].(string),
		}
		mappedStages = append(mappedStages, newStage)
	}
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	if d.HasChange("name") {
		err := c.UpdateApiPipeline(&organizationId, &xc.ApiPipeline{
			Type:           "pipeline",
			ID:             id,
			Name:           name,
			PipelineType:   pipelineType,
			VpcPipelineId:  vpcPipelineId,
			GitRepoId:      gitRepoId,
			Branch:         branch,
			Stages:         mappedStages,
			OrganizationId: organizationId,
			OwningUserId:   owningUserId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceApiPipelineRead(ctx, d, m)
}

func resourceApiPipelineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	err = c.DeleteApiPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceStaticContentPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStaticContentPipelineCreate,
		ReadContext:   resourceStaticContentPipelineRead,
		UpdateContext: resourceStaticContentPipelineUpdate,
		DeleteContext: resourceStaticContentPipelineDelete,
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

func resourceStaticContentPipelineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	cloudProviderId := d.Get("cloud_provider_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	branch := d.Get("branch").(string)
	stages := d.Get("stages").([]interface{})
	mappedStages := []xc.Stage{}
	for _, stage := range stages {
		newStage := xc.Stage {
			Name: stage.(map[string]interface{})["name"].(string),
		}
		mappedStages = append(mappedStages, newStage)
	}
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	location, err := c.CreateStaticContentPipeline(&organizationId, &xc.StaticContentPipeline{
		Type:            "pipeline",
		Name:            name,
		PipelineType:    pipelineType,
		CloudProviderId: cloudProviderId,
		GitRepoId:       gitRepoId,
		Branch:          branch,
		Stages:          mappedStages,
		OrganizationId:  organizationId,
		OwningUserId:    owningUserId,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	index := strings.LastIndex(*location, "/")
	id := string((*location)[(index + 1):])

	d.SetId(id)

	staticContentPipeline, err := c.GetStaticContentPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", staticContentPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", staticContentPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStaticContentPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	staticContentPipeline, err := c.GetStaticContentPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", staticContentPipeline.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", staticContentPipeline.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_type", staticContentPipeline.PipelineType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cloud_provider_id", staticContentPipeline.CloudProviderId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_repo_id", staticContentPipeline.GitRepoId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("branch", staticContentPipeline.Branch); err != nil {
		return diag.FromErr(err)
	}

	stages := make([]interface{}, len(staticContentPipeline.Stages))
	for i, stage := range staticContentPipeline.Stages {
		newStage := make(map[string]interface{})

		newStage["name"] = stage.Name
		stages[i] = newStage
	}
	if err := d.Set("stages", stages); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", staticContentPipeline.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", staticContentPipeline.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", staticContentPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", staticContentPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStaticContentPipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	cloudProviderId := d.Get("cloud_provider_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	branch := d.Get("branch").(string)
	stages := d.Get("stages").([]interface{})
	mappedStages := []xc.Stage{}
	for _, stage := range stages {
		newStage := xc.Stage {
			Name: stage.(map[string]interface{})["name"].(string),
		}
		mappedStages = append(mappedStages, newStage)
	}
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	if d.HasChange("name") {
		err := c.UpdateStaticContentPipeline(&organizationId, &xc.StaticContentPipeline{
			Type:            "pipeline",
			ID:              id,
			Name:            name,
			PipelineType:    pipelineType,
			CloudProviderId: cloudProviderId,
			GitRepoId:       gitRepoId,
			Branch:          branch,
			Stages:          mappedStages,
			OrganizationId:  organizationId,
			OwningUserId:    owningUserId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceStaticContentPipelineRead(ctx, d, m)
}

func resourceStaticContentPipelineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	err := c.DeleteStaticContentPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

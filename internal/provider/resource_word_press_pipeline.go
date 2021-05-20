package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceWordPressPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWordPressPipelineCreate,
		ReadContext:   resourceWordPressPipelineRead,
		UpdateContext: resourceWordPressPipelineUpdate,
		DeleteContext: resourceWordPressPipelineDelete,
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
			"k8s_pipeline_id": {
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

func resourceWordPressPipelineCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	k8sPipelineId := d.Get("k8s_pipeline_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	branch := d.Get("branch").(string)
	stages := d.Get("stages").([]interface{})
	mappedStages := []xc.WordPressStage{}
	for _, stage := range stages {
		newStage := xc.WordPressStage{
			Name: stage.(map[string]interface{})["name"].(string),
		}
		mappedStages = append(mappedStages, newStage)
	}
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	location, err := c.CreateWordPressPipeline(&organizationId, &xc.WordPressPipeline{
		Type:           "pipeline",
		Name:           name,
		PipelineType:   pipelineType,
		K8sPipelineId:  k8sPipelineId,
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

	wordPressPipeline, err := c.GetWordPressPipeline(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", wordPressPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", wordPressPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWordPressPipelineRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	wordPressPipeline, err := c.GetWordPressPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", wordPressPipeline.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", wordPressPipeline.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("pipeline_type", wordPressPipeline.PipelineType); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("k8s_pipeline_id", wordPressPipeline.K8sPipelineId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_repo_id", wordPressPipeline.GitRepoId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("branch", wordPressPipeline.Branch); err != nil {
		return diag.FromErr(err)
	}

	stages := make([]interface{}, len(wordPressPipeline.Stages))
	for i, stage := range wordPressPipeline.Stages {
		newStage := make(map[string]interface{})

		newStage["name"] = stage.Name
		stages[i] = newStage
	}
	if err := d.Set("stages", stages); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", wordPressPipeline.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", wordPressPipeline.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", wordPressPipeline.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", wordPressPipeline.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWordPressPipelineUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	pipelineType := d.Get("pipeline_type").(string)
	k8sPipelineId := d.Get("k8s_pipeline_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	branch := d.Get("branch").(string)
	stages := d.Get("stages").([]interface{})
	mappedStages := []xc.WordPressStage{}
	for _, stage := range stages {
		newStage := xc.WordPressStage{
			Name: stage.(map[string]interface{})["name"].(string),
		}
		mappedStages = append(mappedStages, newStage)
	}
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	if d.HasChange("name") {
		err := c.UpdateWordPressPipeline(&organizationId, &xc.WordPressPipeline{
			Type:           "pipeline",
			ID:             id,
			Name:           name,
			PipelineType:   pipelineType,
			K8sPipelineId:  k8sPipelineId,
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

	return resourceWordPressPipelineRead(ctx, d, m)
}

func resourceWordPressPipelineDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	err := c.DeleteWordPressPipeline(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourcePipelinePrototype() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelinePrototypeCreate,
		ReadContext:   resourcePipelinePrototypeRead,
		UpdateContext: resourcePipelinePrototypeUpdate,
		DeleteContext: resourcePipelinePrototypeDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"references": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameter_definitions": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"terraform": {
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

func resourcePipelinePrototypeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	description := d.Get("description").(string)
	active := d.Get("active").(bool)
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)
	referencesStr := d.Get("references").(string)
	var referencesData []xc.Reference
	err := json.Unmarshal([]byte(referencesStr), &referencesData)
	if err != nil {
		return diag.FromErr(err)
	}
	parameterDefinitionsStr := d.Get("parameter_definitions").(string)
	var parameterDefinitionsData []xc.ParameterDefinition
	err = json.Unmarshal([]byte(parameterDefinitionsStr), &parameterDefinitionsData)
	if err != nil {
		return diag.FromErr(err)
	}
	terraformStr := d.Get("terraform").(string)
	var terraformData xc.Terraform
	err = json.Unmarshal([]byte(terraformStr), &terraformData)
	if err != nil {
		return diag.FromErr(err)
	}

	location, err := c.CreatePipelinePrototype(&organizationId, &xc.PipelinePrototype{
		Type:                 "pipeline-prototype",
		Name:                 name,
		References:           referencesData,
		Version:              version,
		Description:          description,
		Active:               active,
		ParameterDefinitions: parameterDefinitionsData,
		Terraform:            terraformData,
		OrganizationId:       organizationId,
		OwningUserId:         owningUserId,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	id := getIdFromLocationUrl(location)

	d.SetId(*id)

	PipelinePrototype, err := c.GetPipelinePrototype(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", PipelinePrototype.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", PipelinePrototype.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePipelinePrototypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	pipelinePrototype, err := c.GetPipelinePrototype(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", pipelinePrototype.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", pipelinePrototype.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("version", pipelinePrototype.Version); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", pipelinePrototype.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("active", pipelinePrototype.Active); err != nil {
		return diag.FromErr(err)
	}

	referencesData, err := json.Marshal(pipelinePrototype.References)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("references", string(referencesData)); err != nil {
		return diag.FromErr(err)
	}

	parameterDefinitionsData, err := json.Marshal(pipelinePrototype.ParameterDefinitions)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parameter_definitions", string(parameterDefinitionsData)); err != nil {
		return diag.FromErr(err)
	}

	terraformData, err := json.Marshal(pipelinePrototype.Terraform)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("terraform", string(terraformData)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", pipelinePrototype.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", pipelinePrototype.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", pipelinePrototype.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", pipelinePrototype.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePipelinePrototypeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	version := d.Get("version").(string)
	description := d.Get("description").(string)
	active := d.Get("active").(bool)
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)
	referencesStr := d.Get("references").(string)
	var referencesData []xc.Reference
	err := json.Unmarshal([]byte(referencesStr), &referencesData)
	if err != nil {
		return diag.FromErr(err)
	}
	parameterDefinitionsStr := d.Get("parameter_definitions").(string)
	var parameterDefinitionsData []xc.ParameterDefinition
	err = json.Unmarshal([]byte(parameterDefinitionsStr), &parameterDefinitionsData)
	if err != nil {
		return diag.FromErr(err)
	}
	terraformStr := d.Get("terraform").(string)
	var terraformData xc.Terraform
	err = json.Unmarshal([]byte(terraformStr), &terraformData)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		err := c.UpdatePipelinePrototype(&organizationId, &xc.PipelinePrototype{
			Type:                 "pipeline-prototype",
			ID:                   id,
			Name:                 name,
			References:           referencesData,
			Version:              version,
			Description:          description,
			Active:               active,
			ParameterDefinitions: parameterDefinitionsData,
			Terraform:            terraformData,
			OrganizationId:       organizationId,
			OwningUserId:         owningUserId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePipelinePrototypeRead(ctx, d, m)
}

func resourcePipelinePrototypeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	err := c.DeletePipelinePrototype(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

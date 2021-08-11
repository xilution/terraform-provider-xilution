package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourcePipelinePrototype() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePipelinePrototypeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"references": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameter_definitions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"terraform": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owning_user_id": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourcePipelinePrototypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Get("id").(string)

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

	d.SetId(pipelinePrototype.ID)

	return diags
}

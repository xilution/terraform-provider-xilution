package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owning_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
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
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_auth": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_sign_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceOrganizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("id").(string)

	organization, err := c.GetOrganization(&organizationId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", organization.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", organization.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("domain", organization.Domain); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("auth_client_id", organization.AuthClientId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", organization.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", organization.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("active", organization.Active); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", organization.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", organization.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", organization.Url); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("auto_auth", organization.AutoAuth); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("show_sign_up", organization.ShowSignUp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(organization.ID)

	return diags
}

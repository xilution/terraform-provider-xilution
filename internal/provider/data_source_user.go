package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
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
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	userId := d.Get("id").(string)

	user, err := c.GetUser(&organizationId, &userId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", user.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("first_name", user.FirstName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("last_name", user.LastName); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("username", user.Username); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", user.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", user.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("active", user.Active); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", user.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", user.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.ID)

	return diags
}

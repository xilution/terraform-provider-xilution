package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceClient() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClientRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"grants": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"redirect_uris": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"client_user_id": {
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
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	clientId := d.Get("id").(string)

	client, err := c.GetClient(&organizationId, &clientId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", client.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", client.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("grants", client.Grants); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("redirect_uris", client.RedirectUris); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("client_user_id", client.ClientUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", client.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", client.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("active", client.Active); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", client.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", client.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("secret", client.Secret); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(client.ID)

	return diags
}

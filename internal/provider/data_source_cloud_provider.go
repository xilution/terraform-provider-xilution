package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceCloudProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudProviderRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
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

func dataSourceCloudProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	cloudProviderId := d.Get("id").(string)

	cloudProvider, err := c.GetCloudProvider(&organizationId, &cloudProviderId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", cloudProvider.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", cloudProvider.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cloud_provider", cloudProvider.Provider); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("account_id", cloudProvider.AccountId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("region", cloudProvider.Region); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", cloudProvider.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", cloudProvider.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", cloudProvider.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", cloudProvider.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudProvider.ID)

	return diags
}

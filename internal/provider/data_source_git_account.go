package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceGitAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGitAccountRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_provider": {
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

func dataSourceGitAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	gitAccountId := d.Get("id").(string)

	gitAccount, err := c.GetGitAccount(&organizationId, &gitAccountId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", gitAccount.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", gitAccount.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_provider", gitAccount.Provider); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", gitAccount.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", gitAccount.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", gitAccount.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", gitAccount.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(gitAccount.ID)

	return diags
}

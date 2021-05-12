package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceGitRepo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGitRepoRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"git_account_id": {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceGitRepoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	gitAccountId := d.Get("git_account_id").(string)
	gitRepoId := d.Get("id").(string)

	gitRepo, err := c.GetGitRepo(&organizationId, &gitAccountId, &gitRepoId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", gitRepo.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", gitRepo.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_account_id", gitRepo.GitAccountId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", gitRepo.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", gitRepo.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", gitRepo.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", gitRepo.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(gitRepo.ID)

	return diags
}

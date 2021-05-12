package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func dataSourceGitRepoEvent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGitRepoEventRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_repo_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceGitRepoEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	gitAccountId := d.Get("git_account_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	gitRepoEventId := d.Get("id").(string)

	gitRepoEvent, err := c.GetGitRepoEvent(&organizationId, &gitAccountId, &gitRepoId, &gitRepoEventId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", gitRepoEvent.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_account_id", gitRepoEvent.GitAccountId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_repo_id", gitRepoEvent.GitRepoId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", gitRepoEvent.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_type", gitRepoEvent.EventType); err != nil {
		return diag.FromErr(err)
	}

	jsonStr, err := json.Marshal(gitRepoEvent.Parameters)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parameters", string(jsonStr)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("owning_user_id", gitRepoEvent.OwningUserId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", gitRepoEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", gitRepoEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(gitRepoEvent.ID)

	return diags
}

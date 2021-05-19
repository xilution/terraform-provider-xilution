package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceGitRepo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitRepoCreate,
		ReadContext:   resourceGitRepoRead,
		UpdateContext: resourceGitRepoUpdate,
		DeleteContext: resourceGitRepoDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceGitRepoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	organizationId := d.Get("organization_id").(string)
	gitAccountId := d.Get("git_account_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	location, err := c.CreateGitRepo(&organizationId, &xc.GitRepo{
		Type:           "git-repo",
		Name:           name,
		GitAccountId:   gitAccountId,
		OrganizationId: organizationId,
		OwningUserId:   owningUserId,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	index := strings.LastIndex(*location, "/")
	id := string((*location)[(index + 1):])

	d.SetId(id)

	gitRepo, err := c.GetGitRepo(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", gitRepo.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", gitRepo.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceGitRepoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	gitRepo, err := c.GetGitRepo(&organizationId, &id)
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

	return diags
}

func resourceGitRepoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	organizationId := d.Get("organization_id").(string)
	gitAccountId := d.Get("git_account_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	if d.HasChange("name") {
		err := c.UpdateGitRepo(&organizationId, &xc.GitRepo{
			Type:           "git-repo",
			ID:             id,
			Name:           name,
			GitAccountId:   gitAccountId,
			OrganizationId: organizationId,
			OwningUserId:   owningUserId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitRepoRead(ctx, d, m)
}

func resourceGitRepoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	err := c.DeleteGitRepo(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

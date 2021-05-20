package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceGitAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitAccountCreate,
		ReadContext:   resourceGitAccountRead,
		UpdateContext: resourceGitAccountUpdate,
		DeleteContext: resourceGitAccountDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"git_provider": {
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

func resourceGitAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	provider := d.Get("git_provider").(string)
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	location, err := c.CreateGitAccount(&organizationId, &xc.GitAccount{
		Type:           "git-account",
		Name:           name,
		Provider:       provider,
		OrganizationId: organizationId,
		OwningUserId:   owningUserId,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	id := getIdFromLocationUrl(location)

	d.SetId(*id)

	gitAccount, err := c.GetGitAccount(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", gitAccount.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", gitAccount.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceGitAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	gitAccount, err := c.GetGitAccount(&organizationId, &id)
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

	return diags
}

func resourceGitAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	provider := d.Get("git_provider").(string)
	organizationId := d.Get("organization_id").(string)
	owningUserId := d.Get("owning_user_id").(string)

	if d.HasChange("name") {
		err := c.UpdateGitAccount(&organizationId, &xc.GitAccount{
			Type:           "git-account",
			ID:             id,
			Name:           name,
			Provider:       provider,
			OrganizationId: organizationId,
			OwningUserId:   owningUserId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitAccountRead(ctx, d, m)
}

func resourceGitAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	err := c.DeleteGitAccount(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

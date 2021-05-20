package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceCloudProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProviderCreate,
		ReadContext:   resourceCloudProviderRead,
		UpdateContext: resourceCloudProviderUpdate,
		DeleteContext: resourceCloudProviderDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
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

func resourceCloudProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	organizationId := d.Get("organization_id").(string)
	provider := d.Get("cloud_provider").(string)
	accountId := d.Get("account_id").(string)
	region := d.Get("region").(string)
	owningUserId := d.Get("owning_user_id").(string)

	location, err := c.CreateCloudProvider(&organizationId, &xc.CloudProvider{
		Type:           "cloud-provider",
		Name:           name,
		Provider:       provider,
		AccountId:      accountId,
		Region:         region,
		OrganizationId: organizationId,
		OwningUserId:   owningUserId,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	id := getIdFromLocationUrl(location)

	d.SetId(*id)

	cloudProvider, err := c.GetCloudProvider(&organizationId, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", cloudProvider.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", cloudProvider.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCloudProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	cloudProvider, err := c.GetCloudProvider(&organizationId, &id)
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

	return diags
}

func resourceCloudProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	id := d.Id()
	name := d.Get("name").(string)
	organizationId := d.Get("organization_id").(string)
	provider := d.Get("cloud_provider").(string)
	accountId := d.Get("account_id").(string)
	region := d.Get("region").(string)
	owningUserId := d.Get("owning_user_id").(string)

	if d.HasChange("name") {
		err := c.UpdateCloudProvider(&organizationId, &xc.CloudProvider{
			Type:           "cloud-provider",
			ID:             id,
			Name:           name,
			Provider:       provider,
			AccountId:      accountId,
			Region:         region,
			OrganizationId: organizationId,
			OwningUserId:   owningUserId,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCloudProviderRead(ctx, d, m)
}

func resourceCloudProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	id := d.Id()

	err := c.DeleteCloudProvider(&organizationId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

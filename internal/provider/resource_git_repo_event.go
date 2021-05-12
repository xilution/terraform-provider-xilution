package provider

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	xc "github.com/xilution/xilution-client-go"
)

func resourceGitRepoEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitRepoEventCreate,
		ReadContext:   resourceGitRepoEventRead,
		UpdateContext: resourceGitRepoEventUpdate,
		DeleteContext: resourceGitRepoEventDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
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
				Required: true,
			},
			"parameters": {
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

func resourceGitRepoEventCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	gitAccountId := d.Get("git_account_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	owningUserId := d.Get("owning_user_id").(string)
	parametersStr := d.Get("parameters").(string)
	var parametersData map[string]interface{}
	err := json.Unmarshal([]byte(parametersStr), &parametersData)
	if err != nil {
		return diag.FromErr(err)
	}

	eventType := d.Get("event_type").(string)

	location, err := c.CreateGitRepoEvent(&organizationId, &gitAccountId, &gitRepoId, &xc.GitRepoEvent{
		Type:           "git-repo-event",
		GitAccountId:   gitAccountId,
		GitRepoId:      gitRepoId,
		OrganizationId: organizationId,
		OwningUserId:   owningUserId,
		Parameters:     parametersData,
		EventType:      eventType,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	index := strings.LastIndex(*location, "/")
	id := string((*location)[(index + 1):])

	d.SetId(id)

	gitRepoEvent, err := c.GetGitRepoEvent(&organizationId, &gitAccountId, &gitRepoId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("created_at", gitRepoEvent.CreatedAt); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("modified_at", gitRepoEvent.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceGitRepoEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*xc.XilutionClient)

	var diags diag.Diagnostics

	organizationId := d.Get("organization_id").(string)
	gitAccountId := d.Get("git_account_id").(string)
	gitRepoId := d.Get("git_repo_id").(string)
	id := d.Id()

	gitRepoEvent, err := c.GetGitRepoEvent(&organizationId, &gitAccountId, &gitRepoId, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", gitRepoEvent.ID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("git_account_id", gitRepoEvent.GitAccountId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_id", gitRepoEvent.OrganizationId); err != nil {
		return diag.FromErr(err)
	}

	jsonStr, err := json.Marshal(gitRepoEvent.Parameters)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Println("[DEBUG] jsonStr: " + string(jsonStr))
	if err := d.Set("parameters", string(jsonStr)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("event_type", gitRepoEvent.EventType); err != nil {
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

	return diags
}

func resourceGitRepoEventUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceGitRepoEventRead(ctx, d, m)
}

func resourceGitRepoEventDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

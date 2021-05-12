terraform {
  required_providers {
    xilution = {
      version = "0.1"
      source  = "xilution.com/general/xilution"
    }
  }
}

locals {
  organization_id = "1645137c2f53427da00edaa680256215"
  client_id       = "40808d61f755421098b64529af53ceb9"
  user_id         = "9bfbcd5e1baf490694170cd623e71305"
}

provider "xilution" {
  organization_id = local.organization_id
  client_id       = local.client_id
}

data "xilution_organization" "xilution" {
  id = local.organization_id
}

output "xilution_organization" {
  value = data.xilution_organization.xilution
}

data "xilution_client" "terraform_client" {
  id              = local.client_id
  organization_id = local.organization_id
}

output "xilution_client" {
  value = data.xilution_client.terraform_client
}

data "xilution_user" "tbrunia" {
  id              = local.user_id
  organization_id = local.organization_id
}

output "xilution_user" {
  value = data.xilution_user.tbrunia
}

resource "xilution_git_account" "xilution_git_account" {
  name            = "xilution"
  git_provider    = "GIT_HUB"
  organization_id = local.organization_id
  owning_user_id  = local.user_id
}

data "xilution_git_account" "xilution_git_account" {
  id = xilution_git_account.xilution_git_account.id
  organization_id = local.organization_id
}

output "xilution_git_account" {
  value = data.xilution_git_account.xilution_git_account
}

resource "xilution_git_repo" "xilution_temp_git_repo" {
  name            = "xilution-temp"
  organization_id = local.organization_id
  owning_user_id  = local.user_id
  git_account_id  = xilution_git_account.xilution_git_account.id
}

data "xilution_git_repo" "xilution_temp_git_repo" {
  id = xilution_git_repo.xilution_temp_git_repo.id
  organization_id = local.organization_id
  git_account_id  = xilution_git_account.xilution_git_account.id
}

output "xilution_git_repo" {
  value = data.xilution_git_repo.xilution_temp_git_repo
}

resource "xilution_git_repo_event" "xilution_temp_git_repo_event" {
  organization_id = local.organization_id
  owning_user_id  = local.user_id
  git_account_id  = xilution_git_account.xilution_git_account.id
  git_repo_id     = xilution_git_repo.xilution_temp_git_repo.id
  event_type      = "CREATE_REPO_FROM_TEMPLATE_REPO"
  parameters      = jsonencode({
    "sourceOwner": "xilution",
    "sourceRepo": "xilution-bison-poc-template",
    "description": "A new repo from a copy",
    "commitMessage": "Initial repo setup",
    "isPrivate": true,
    "params": "{\"world\": \"planet\"}"
  })
}

data "xilution_git_repo_event" "xilution_temp_git_repo_event" {
  id = xilution_git_repo_event.xilution_temp_git_repo_event.id
  organization_id = local.organization_id
  git_account_id  = xilution_git_account.xilution_git_account.id
  git_repo_id     = xilution_git_repo.xilution_temp_git_repo.id
}

output "xilution_temp_git_repo_event" {
  value = data.xilution_git_repo_event.xilution_temp_git_repo_event
}

resource "xilution_cloud_provider" "xilution_aws_prod" {
  organization_id = local.organization_id
  owning_user_id  = local.user_id
  name = "Xilution AWS (Prod)"
  cloud_provider = "AWS"
  account_id = "952573012699"
  region = "us-east-1"
}

data "xilution_cloud_provider" "xilution_aws_prod" {
  id = xilution_cloud_provider.xilution_aws_prod.id
  organization_id = local.organization_id
}

output "xilution_aws_prod" {
  value = data.xilution_cloud_provider.xilution_aws_prod
}

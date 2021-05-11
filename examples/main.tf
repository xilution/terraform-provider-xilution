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

output "xilution_git_account" {
  value = xilution_git_account.xilution_git_account
}

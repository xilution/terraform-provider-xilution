terraform {
  required_providers {
    xilution = {
      # source = "xilution.com/xilution/xilution"
      # version = "0.1.6"
      source  = "registry.terraform.io/xilution/xilution"
      version = "0.1.6"
    }
  }
}

locals {
  organization_id = var.ORGANIZATION_ID
  client_id       = var.CLIENT_ID
  client_secret   = var.CLIENT_SECRET
  user_id         = var.USER_ID
}

provider "xilution" {
  organization_id = local.organization_id
  client_id       = local.client_id
  client_secret   = local.client_secret
}

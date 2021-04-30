terraform {
  required_providers {
    xilution = {
      version = "0.1"
      source  = "xilution.com/general/xilution"
    }
  }
}

provider "xilution" {
  username = "education"
  password = "test123"
}

module "psl" {
  source = "./coffee"

  coffee_name = "Packer Spiced Latte"
}

output "psl" {
  value = module.psl.coffee
}

data "xilution_order" "order" {
  id = 1
}

output "order" {
  value = data.xilution_order.order
}

terraform {
  required_providers {
    xilution = {
      version = "0.1"
      source  = "xilution.com/general/xilution"
    }
  }
}

variable "coffee_name" {
  type    = string
  default = "Vagrante espresso"
}

data "xilution_coffees" "all" {}

# Returns all coffees
output "all_coffees" {
  value = data.xilution_coffees.all.coffees
}

# Only returns packer spiced latte
output "coffee" {
  value = {
    for coffee in data.xilution_coffees.all.coffees :
    coffee.id => coffee
    if coffee.name == var.coffee_name
  }
}

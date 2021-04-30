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

resource "xilution_order" "sample" {}

output "sample_order" {
  value = xilution_order.sample
}

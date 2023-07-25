#vpc
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "zone-name"
    values = ["ap-northeast-2a", "ap-northeast-2c"]
  }
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "4.0.1"

  name = var.vpc_name
  cidr = var.vpc_cidr

  azs                  = data.aws_availability_zones.available.names
  private_subnets      = var.vpc_privatesubnet
  public_subnets       = var.vpc_publicsubnet
  enable_dns_hostnames = true
  enable_dns_support   = true

  enable_nat_gateway = true
  single_nat_gateway = false
}
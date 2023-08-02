#vpc
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "4.0.1"

  name = var.vpc_name
  cidr = var.vpc_cidr

  azs                  = var.vpc_azs
  private_subnets      = var.vpc_privatesubnet
  public_subnets       = var.vpc_publicsubnet
  enable_dns_hostnames = true
  enable_dns_support   = true

  enable_nat_gateway = true
  single_nat_gateway = false

  manage_default_security_group = false
}
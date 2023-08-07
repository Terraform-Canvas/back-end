#vpc
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "4.0.1"

  name = var.vpc_name
  cidr = var.vpc_cidr

  azs                  = var.vpc_azs
  public_subnets       = var.vpc_publicsubnet
  private_subnets      = var.vpc_privatesubnet
  enable_dns_hostnames = true
  enable_dns_support   = true

  enable_nat_gateway = length(var.vpc_publicsubnet) == 0 ? false : true
  single_nat_gateway = false
}
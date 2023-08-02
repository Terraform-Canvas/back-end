
module "eks_control" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~>19.0"

  cluster_name    = var.eks_control_name
  cluster_version = var.eks_control_version

  cluster_endpoint_public_access = true

  cluster_addons = {
    coredns    = {}
    kube-proxy = {}
    vpc-cni    = {}
  }

  vpc_id                   = module.vpc.vpc_id
  control_plane_subnet_ids = var.eks_control_subnet_type == "privatesubnet" ? [for i in range(var.eks_control_subnet_count[0], var.eks_control_subnet_count[1]) : module.vpc.private_subnets[i]] : [for i in range(var.eks_control_subnet_count[0], var.eks_control_subnet_count[1]) : module.vpc.public_subnets[i]]

}

module "eks_worker" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~>19.0"

  cluster_name    = var.eks_control_name
  cluster_version = var.eks_control_version

  create_cloudwatch_log_group = false
  create_kms_key              = false

  cluster_endpoint_public_access  = true
  cluster_endpoint_private_access = true

  vpc_id     = module.vpc2.vpc_id
  subnet_ids = var.eks_worker_subnet_type == "privatesubnet" ? [for i in range(var.eks_worker_subnet_count[0], var.eks_worker_subnet_count[1]) : module.vpc.private_subnets[i]] : [for i in range(var.eks_worker_subnet_count[0], var.eks_worker_subnet_count[1]) : module.vpc.public_subnets[i]]

  #workers.tf
}
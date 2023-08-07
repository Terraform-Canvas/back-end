#eks
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~>19.0"

  cluster_name    = var.eks_name
  cluster_version = var.eks_version

  cluster_endpoint_public_access  = var.eks_endpoint_public
  cluster_endpoint_private_access = var.eks_endpoint_private
  create_cloudwatch_log_group     = false

  cluster_addons = {
    coredns = {
      resolve_conflicts = "OVERWRITE"
    }
    kube-proxy = {}
    vpc-cni = {
      resolve_conflicts = "OVERWRITE"
    }
  }

  vpc_id     = module.vpc.vpc_id
  subnet_ids = var.eks_subnet_type == "privatesubnet" ? [for i in range(var.eks_subnet_count[0], var.eks_subnet_count[1]) : module.vpc.private_subnets[i]] : [for i in range(var.eks_subnet_count[0], var.eks_subnet_count[1]) : module.vpc.public_subnets[i]]

  //auth
  manage_aws_auth_configmap = true
  aws_auth_users = [
    {
      userarn  = var.eks_userarn
      username = var.eks_username
    },
  ]

  //node_groups
}

data "aws_eks_cluster_auth" "this" {
  name = module.eks.cluster_name
}

provider "kubernetes" {
  host                   = module.eks.cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks.cluster_certificate_authority_data)
  token                  = data.aws_eks_cluster_auth.this.token
}
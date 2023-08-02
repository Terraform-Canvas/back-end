  #eks_managed_main
  /*
  eks_managed_node_groups = {
    blue = {}
    green = {
      min_size     = var.eks_managed_min_size
      max_size     = var.eks_managed_max_size
      desired_size = var.eks_managed_desired_size

      instance_types = var.eks_managed_instance_types
      disk_size      = var.eks_managed_disk_size
      capacity_type  = var.eks_managed_capacity_type
    }
  }
  */

  #eks_managed_variable
  /*
  variable "eks_managed_min_size" {
    type = number
  }

  variable "eks_managed_max_size" {
    type = number
  }

  variable "eks_managed_desired_size" {
    type = number
  }

  variable "eks_managed_instance_types" {
    type = list(any)
  }

  variable "eks_managed_disk_size" {
    type = number
  }

  variable "eks_managed_capacity_type" {
    type = string
  }

  */

  #fargate_main
  /*
  fargate_profiles = {
    default = {
      name = var.fargate_name
      selectors = [
        for ns in var.fargate_namespaces : {
          namespace = ns
        }
      ]
    }
  }
  */

  #fargate_variable
  /*
  variable "fargate_name" {
    type = string
  }

  variable "fargate_namespaces" {
    type = list(any)
  }
  */
variable "eks_name" {
  type = string
}

variable "eks_version" {
  type = string
}

variable "eks_subnet_count" {
  type = list(any)
}

variable "eks_subnet_type" {
  type = string
}

variable "eks_userarn" {
  type = string
}

variable "eks_username" {
  type = string
}

variable "eks_endpoint_private" {
  type = bool
  default = true
}

variable "eks_endpoint_public" {
  type = bool
  default = true
}
variable "eks_control_name" {
  type = string
}

variable "eks_control_version" {
  type = string
}

variable "eks_control_subnet_count" {
  type = list(any)
}

variable "eks_worker_subnet_count" {
  type = list(any)
}
variable "vpc_name" {
  type = string
}

variable "vpc_azs" {
  type = list(any)
}

variable "vpc_cidr" {
  type = string
}

variable "vpc_publicsubnet" {
  type = list(any)
}

variable "vpc_privatesubnet" {
  type = list(any)
}
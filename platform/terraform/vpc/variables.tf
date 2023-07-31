variable "vpc_name" {
  type    = string
  default = "example-vpc"
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "vpc_publicsubnet" {
  type = list(any)
}

variable "vpc_privatesubnet" {
  type = list(any)
}

variable "vpc_azs" {
  type    = list(any)
  default = ["us-east-1a", "us-east-1c"]
}

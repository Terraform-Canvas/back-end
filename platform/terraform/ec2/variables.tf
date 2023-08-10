variable "ec2_name" {
  type = string
}

variable "ec2_type" {
  type = string
}

variable "ec2_key" {
  type = string
}

variable "ec2_public_ip" {
  type = bool
}

variable "ec2_ami" {
  type = string
}

variable "ec2_subnet_type" {
    type = string
}

variable "ec2_subnet_count" {
    type = list(any)
}
variable "asg_min_size" {
  type    = number
  default = 1
}

variable "asg_max_size" {
  type    = number
  default = 3
}

variable "asg_instance_type" {
  type    = string
  default = "t2.micro"
}

variable "asg_image_id" {
  type = string
}

variable "asg_subnet_count" {
  type    = list(number)
  default = [0, 2]
}

variable "asg_subnet_type" {
  type    = string
  default = "privatesubnet"
}

variable "asg_desired_capacity" {
  type    = number
  default = 2
}
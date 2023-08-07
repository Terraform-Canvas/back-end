variable "alb_subnet_count" {
  type    = list(number)
  default = [0, 2]
}

variable "alb_subnet_type" {
  type    = string
  default = "publicsubnet"
}
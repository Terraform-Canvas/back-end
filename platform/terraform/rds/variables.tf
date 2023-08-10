variable "rds_name" {
  type = string
}

variable "rds_engine" {
  type = string
}

variable "rds_instance_class" {
  type = string
}

variable "rds_engine_version" {
  type = string
}

variable "rds_allocated_storage" {
  type = number
}

variable "rds_username" {
  type = string
}

variable "rds_port" {
  type = string
}

variable "rds_public_access" {
  type = bool
}

variable "rds_multi_az" {
  type = bool
}

variable "rds_subnet_count" {
  type = list(any)
}

variable "rds_subnet_type" {
  type = string
}
#rds
module "rds" {
  source     = "terraform-aws-modules/rds/aws"
  identifier = var.rds_name

  engine                    = var.rds_engine
  engine_version            = var.rds_engine_version
  instance_class            = var.rds_instance_class
  allocated_storage         = var.rds_allocated_storage
  create_db_parameter_group = false
  create_db_option_group    = false

  db_name             = var.rds_name
  username            = var.rds_username
  port                = var.rds_port
  publicly_accessible = var.rds_public_access

  create_db_subnet_group = true
  subnet_ids             = var.rds_subnet_type == "privatesubnet" ? [for i in range(var.rds_subnet_count[0], var.rds_subnet_count[1]) : module.vpc.private_subnets[i]] : [for i in range(var.rds_subnet_count[0], var.rds_subnet_count[1]) : module.vpc.public_subnets[i]]

  multi_az               = var.rds_multi_az
  vpc_security_group_ids = [module.rds_sg.security_group_id]
}
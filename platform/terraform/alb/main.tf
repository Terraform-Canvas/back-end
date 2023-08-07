#alb 
module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 8.0"

  name = "tc-alb"

  load_balancer_type = "application"

  vpc_id                = module.vpc.vpc_id
  subnets               = var.alb_subnet_type == "privatesubnet" ? [for i in range(var.asg_subnet_count[0], var.asg_subnet_count[1]) : module.vpc.private_subnets[i]] : [for i in range(var.asg_subnet_count[0], var.asg_subnet_count[1]) : module.vpc.public_subnets[i]]
  security_groups       = [module.alb_sg.security_group_id]
  internal              = false
  create_security_group = false

  http_tcp_listeners = [
    {
      port               = 80
      protocol           = "HTTP"
      target_group_index = 0
    }
  ]

  target_groups = [{
    name             = "Terraform-Canvas-ALB-http"
    backend_port     = 80
    backend_protocol = "HTTP"
    target_type      = "instance"
  }]
}
#alb 
module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 8.0"

  name = "tc-alb"

  load_balancer_type = "application"

  vpc_id          = module.vpc.vpc_id
  subnets         = var.alb_subnet
  security_groups = [module.alb_sg.security_group_id]
  internal        = false

  http_tcp_listeners = [
    {
      port        = 80
      protocol    = "HTTP"
      action_type = "forward"
    }
  ]

  target_groups = [{
    name             = "Terraform-Canvas ALB"
    backend_port     = 80
    backend_protocol = "HTTP"
  }]
}
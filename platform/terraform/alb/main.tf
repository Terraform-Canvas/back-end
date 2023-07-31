#alb 
module "alb" {
  source  = "terraform-aws-modules/alb/aws"
  version = "~> 8.0"

  name = "tc-alb"

  load_balancer_type = "application"

  vpc_id          = module.vpc.vpc_id
  subnets         = var.alb_subnet_type == "privatesubnet" ? [for i in range(var.asg_subnet_count[0], var.asg_subnet_count[1]) : module.vpc.private_subnets[i]] : [for i in range(var.asg_subnet_count[0], var.asg_subnet_count[1]) : module.vpc.public_subnets[i]]
  security_groups = [aws_security_group.alb_sg.id]
  internal        = false

  http_tcp_listeners = [
    {
      port        = 80
      protocol    = "HTTP"
      action_type = "forward"
    }
  ]

  target_groups = [{
    name             = "Terraform-Canvas-ALB"
    backend_port     = 80
    backend_protocol = "HTTP"

    health_check = {
      path                = "/"
      protocol            = "HTTP"
      matcher             = "200"
      interval            = 120
      timeout             = 60
      healthy_threshold   = 2
      unhealthy_threshold = 2
    }
  }]
}
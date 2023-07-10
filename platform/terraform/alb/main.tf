#alb 
module "alb" {
    source = "terraform-aws-modules/alb/aws"
    version = "~> 8.0"

    name = "terramino"
    
    load_balancer_type = "application"
    
    vpc_id = module.vpc.vpd_id
    subnets = module.vpc.public_subnets
    security_groups = [module.terramino_lb.id]
    internal = false

    http_tcp_listeners=[
        {
            port = 80
            protocol="HTTP"
            action_type = "forward"
        }
    ]

    target_groups = [
        name = "learn-asg-terramino"
        backend_port = 80
        backend_protocol = "HTTP"
    ]
}
resource "aws_autoscaling_attachment" "terramino" {
  autoscaling_group_name = aws_autoscaling_group.terramino.id
  alb_target_group_arn   = aws_lb_target_group.terramino.arn
}

module "terramino_instance" {
    source = "terraform-aws-modules/security-group/aws//modules/http-80"
    
    name = "learn-asg-terramino-instance"
    vpc_id = module.vpc.vpc_id

    ingress_with_security_group_id = module.terramino_lb.id
}
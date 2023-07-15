# asg
module "asg" {
  source = "terraform-aws-modules/autoscaling/aws"

  name                = "terramino"
  min_size            = var.asg_min_size
  max_size            = var.asg_max_size
  desired_capacity    = var.asg_desired_capacity
  vpc_zone_identifier = module.vpc.private_subnets

  launch_template_name = "learn-terraform-aws-asg-"
  use_name_prefix      = true
  image_id             = data.aws_ami.amazon-linux.id
  instance_type        = var.asg_instance_type
  user_data            = file("user-data.sh")
  security_groups      = [module.terramino_instance.security_group_id]

  tags = {
    key                 = "Name"
    value               = "HashiCorp Learn ASG - Terramino"
    propagate_at_launch = true
  }
}

module "terramino_lb" {
  source = "terraform-aws-modules/security-group/aws//modules/http-80"

  name   = "learn-asg-terramino-lb"
  vpc_id = module.vpc.vpc_id

  ingress_cidr_blocks = ["0.0.0.0/0"]
}
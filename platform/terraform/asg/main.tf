# asg
module "asg" {
  source = "terraform-aws-modules/autoscaling/aws"

  name                = "tc-asg"
  min_size            = var.asg_min_size
  max_size            = var.asg_max_size
  desired_capacity    = var.asg_desired_capacity
  vpc_zone_identifier = module.vpc.private_subnets

  launch_template_name = "learn-terraform-aws-asg-"
  use_name_prefix      = true
  image_id             = var.asg_image_id
  instance_type        = var.asg_instance_type
  user_data            = file("user-data.sh")
  security_groups      = [module.asg_ssg.security_group_id]

  tags = {
    key                 = "Name"
    value               = "Terraform-Canvas ASG"
    propagate_at_launch = true
  }
}

resource "aws_autoscaling_attachment" "asg_alb" {
  autoscaling_group_name = module.asg.autoscaling_group_id
  lb_target_group_arn    = module.alb.target_group_arns[0]
}
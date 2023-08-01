# asg
module "asg" {
  source = "terraform-aws-modules/autoscaling/aws"

  name                      = "tc-asg"
  min_size                  = var.asg_min_size
  max_size                  = var.asg_max_size
  desired_capacity          = var.asg_desired_capacity
  health_check_type         = "ELB"
  health_check_grace_period = 300
  vpc_zone_identifier       = var.asg_subnet_type == "privatesubnet" ? [for i in range(var.asg_subnet_count[0], var.asg_subnet_count[1]) : module.vpc.private_subnets[i]] : [for i in range(var.asg_subnet_count[0], var.asg_subnet_count[1]) : module.vpc.public_subnets[i]]

  launch_template_name = "learn-terraform-aws-asg-"
  use_name_prefix      = true
  image_id             = var.asg_image_id
  instance_type        = var.asg_instance_type
  user_data            = base64encode(file("user-data.sh"))
  security_groups      = [aws_security_group.asg_ssg.id]

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

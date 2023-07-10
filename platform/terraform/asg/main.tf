# asg
module "asg" {
    source = "terraform-aws-modules/autoscaling/aws"

    name = "terramino"
    min_size = var.asg_min_size ####*
    max_size = var.asg_max_size ####*
    desired_capacity = var.asg_desired_capacity ####*
    vpc_zone_identifier = module.vpc.private_subnets

    launch_template_name = "learn-terraform-aws-asg-" #이건 그냥 asg이름-username-asg-로 하는게 어지
    use_name_prefix = true
    image_id = var.asg_image_id
    instance_type = var.asg_instance_type        ####*
    user_data = file("user-data.sh")  ####*
    security_groups = [module.terramino_lb.id]

    tags = {
        key = "Name"
        value = "HashiCorp Learn ASG - Terramino"
        propagate_at_launch = true
    }
}

module "terramino_lb" {
  source = "terraform-aws-modules/security-group/aws//modules/http-80"

  name        = "learn-asg-terramino-lb"
  vpc_id      = module.vpc.vpd_id

  ingress_cidr_blocks = ["10.10.0.0/16"]
}

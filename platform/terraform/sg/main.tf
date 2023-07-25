module "asg_ssg" {
  source = "terraform-aws-modules/security-group/aws"
  name   = "asg_soucesg_sg"
  vpc_id = module.vpc.vpc_id

  computed_ingress_with_source_security_group_id = [
    {
      rule                     = "http-80-tcp"
      source_security_group_id = module.alb_sg.security_group_id
    }
  ]
  number_of_computed_ingress_with_source_security_group_id = 1
}

module "alb_sg" {
  source = "terraform-aws-modules/security-group/aws//modules/http-80"

  name   = "alb_cidr_sg"
  vpc_id = module.vpc.vpc_id

  ingress_cidr_blocks = ["0.0.0.0/0"]
}
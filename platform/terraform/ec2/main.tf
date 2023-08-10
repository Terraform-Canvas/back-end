#ec2
module "ec2" {
  source = "terraform-aws-modules/ec2-instance/aws"

  name = var.ec2_name

  instance_type          = var.ec2_type
  key_name               = var.ec2_key
  associate_public_ip_address = var.ec2_public_ip
  ami = var.ec2_ami
  user_data_base64 = base64encode(file("user-data.sh"))

  vpc_security_group_ids = [module.ec2_sg.security_group_id]
  subnet_id              = var.ec2_subnet_type == "privatesubnet" ? [for i in range(var.ec2_subnet_count[0], var.ec2_subnet_count[1]) : module.vpc.private_subnets[i]][0] : [for i in range(var.ec2_subnet_count[0], var.ec2_subnet_count[1]) : module.vpc.public_subnets[i]][0]
}
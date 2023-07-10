module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "2.77.0"

  name = var.vpc_name      ####*
  cidr = var.vpc_cidr

  azs                  = data.aws_availability_zones.available.names
  private_subnets = var.vpc_privatesubnet
  public_subnets  = var.vpc_publicsubnet
  enable_dns_hostnames = true # VPC 내의 ec2 인스터스에 대해 자동으로 DNS 호스트 이름 할당
  enable_dns_support   = true # 인스턴스는 아마존 내부 DNS 서비스를 통해 도메인 이름 해석 가능 

  #nat를 위한 설정(자동으로 public에만 들어감)
	## 나트 수에 따라 달라짐 -> 표시 방식을 먼저 정하기(간소화 or 아이콘 사용) 
  enable_nat_gateway = true
  single_nat_gateway = false 
}

#ami는 만드는거까진 하지 말자 
data "aws_ami" "amazon-linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-hvm-*-x86_64-ebs"]
  }
}

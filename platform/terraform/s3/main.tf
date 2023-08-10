#s3
module "s3" {
  source = "terraform-aws-modules/s3-bucket/aws"

  bucket                  = var.s3_name
  block_public_acls       = var.s3_public_access
  block_public_policy     = var.s3_public_access
  ignore_public_acls      = var.s3_public_access
  restrict_public_buckets = var.s3_public_access

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  versioning = {
    enabled = true
  }
}

provider "aws" {
  region = var.aws_region

  access_key = "048408301323"
  secret_key = "test"

  endpoints {
    ec2      = "http://localhost:4566"
    rds      = "http://localhost:4566"
    dynamodb = "http://localhost:4566"
    sqs      = "http://localhost:4566"
  }

  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}

module "networking" {
  source = "./modules/networking"

  project            = var.project
  environment        = var.environment
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
}

module "rds" {
  source = "./modules/rds"

  project     = var.project
  environment = var.environment

  # Conexão com a VPC
  vpc_id            = module.networking.vpc_id
  vpc_cidr          = var.vpc_cidr
  private_subnet_ids = module.networking.private_subnet_ids

  # Banco
  database_name   = var.database_name
  master_username = var.master_username
  master_password = var.master_password

  # Configurações por ambiente (use valores diferentes para dev/staging/prod)
  engine_version        = var.engine_version
  engine_mode           = var.engine_mode
  instance_class        = var.instance_class
  instance_count        = var.instance_count
  min_capacity          = var.min_capacity
  max_capacity          = var.max_capacity
  backup_retention_days = var.backup_retention_days
  deletion_protection   = var.deletion_protection
  apply_immediately     = var.apply_immediately
}

module "dynamodb" {
  source = "./modules/dynamodb"

  project     = var.project
  environment = var.environment

  deletion_protection = var.dynamodb_deletion_protection
  ttl_enabled         = var.dynamodb_ttl_enabled
  enable_audit        = var.dynamodb_enable_audit
}

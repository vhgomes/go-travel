# Security Group para o RDS (permite apenas tráfego interno da VPC)
resource "aws_security_group" "rds" {
  name        = "${var.project}-${var.environment}-rds-sg"
  description = "Security group para o RDS Aurora"
  vpc_id      = var.vpc_id

  ingress {
    description     = "Permitir acesso ao PostgreSQL vindo da VPC"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    cidr_blocks     = [var.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-${var.environment}-rds-sg"
  }
}

# Subnet group (associa as subnets privadas ao RDS)
resource "aws_db_subnet_group" "this" {
  name       = "${var.project}-${var.environment}-db-subnet"
  subnet_ids = var.private_subnet_ids

  tags = {
    Name = "${var.project}-${var.environment}-db-subnet"
  }
}

# Cluster Aurora PostgreSQL (Serverless v2 ou provisionado)
resource "aws_rds_cluster" "this" {
  cluster_identifier = "${var.project}-${var.environment}-cluster"
  engine             = "aurora-postgresql"
  engine_version     = var.engine_version
  database_name      = var.database_name
  master_username    = var.master_username
  master_password    = var.master_password

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.this.name

  # Backup e retenção
  backup_retention_period = var.backup_retention_days
  preferred_backup_window = "03:00-04:00"
  preferred_maintenance_window = "sun:05:00-sun:06:00"

  # Deletion protection em produção (opcional via variável)
  deletion_protection = var.deletion_protection

  # Escalabilidade (Serverless v2 ou provisionada)
  serverlessv2_scaling_configuration {
    min_capacity = var.min_capacity
    max_capacity = var.max_capacity
  }

  # Engine mode (provisioned ou serverless)
  engine_mode = var.engine_mode

  # Parâmetros opcionais
  apply_immediately = var.apply_immediately

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}

# Instância (pode ter 1 ou mais)
resource "aws_rds_cluster_instance" "this" {
  count = var.instance_count

  identifier         = "${var.project}-${var.environment}-instance-${count.index + 1}"
  cluster_identifier = aws_rds_cluster.this.id
  instance_class     = var.instance_class
  engine             = aws_rds_cluster.this.engine
  engine_version     = aws_rds_cluster.this.engine_version

  publicly_accessible = false # NUNCA expor o banco à internet

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}

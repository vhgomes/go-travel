variable "aws_region" {
  description = "Região AWS"
  type        = string
  default     = "us-east-1"
}

variable "project" {
  description = "Nome do projeto"
  type        = string
  default     = "travelgo"
}

variable "environment" {
  description = "Ambiente (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "vpc_cidr" {
  description = "CIDR da VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Zonas de disponibilidade"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}

# RDS Variables
variable "database_name" {
  description = "Nome do banco de dados"
  type        = string
  default     = "travelgo"
}

variable "master_username" {
  description = "Usuário mestre do banco"
  type        = string
  default     = "admin"
}

variable "master_password" {
  description = "Senha mestre do banco (use variável de ambiente TF_VAR_master_password)"
  type        = string
  sensitive   = true
  default     = "admin123"  # Mude isso!
}

variable "engine_version" {
  description = "Versão do Aurora PostgreSQL"
  type        = string
  default     = "16.4"
}

variable "engine_mode" {
  description = "Modo do engine: 'provisioned' ou 'serverless'"
  type        = string
  default     = "serverless"
}

variable "instance_class" {
  description = "Classe da instância"
  type        = string
  default     = "db.t4g.small"
}

variable "instance_count" {
  description = "Número de instâncias no cluster"
  type        = number
  default     = 1
}

variable "min_capacity" {
  description = "Capacidade mínima (Serverless v2)"
  type        = number
  default     = 0.5
}

variable "max_capacity" {
  description = "Capacidade máxima (Serverless v2)"
  type        = number
  default     = 4
}

variable "backup_retention_days" {
  description = "Dias de retenção de backup"
  type        = number
  default     = 7
}

variable "deletion_protection" {
  description = "Proteger contra deleção acidental"
  type        = bool
  default     = false
}

variable "apply_immediately" {
  description = "Aplicar mudanças imediatamente"
  type        = bool
  default     = false
}

# DynamoDB Variables
variable "dynamodb_deletion_protection" {
  description = "Proteger tabelas DynamoDB contra deleção acidental"
  type        = bool
  default     = false
}

variable "dynamodb_ttl_enabled" {
  description = "Habilitar TTL na tabela principal de estado"
  type        = bool
  default     = true
}

variable "dynamodb_enable_audit" {
  description = "Criar tabela de auditoria para rastrear mudanças de estado"
  type        = bool
  default     = false
}

# SQS Variables
variable "sqs_max_receive_count" {
  description = "Número máximo de tentativas antes de enviar para DLQ"
  type        = number
  default     = 3
}

variable "sqs_enable_notifications" {
  description = "Criar fila de notificações (não FIFO)"
  type        = bool
  default     = true
}

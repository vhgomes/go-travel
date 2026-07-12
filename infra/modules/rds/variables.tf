variable "project" {
  description = "Nome do projeto"
  type        = string
}

variable "environment" {
  description = "Ambiente (dev, staging, prod)"
  type        = string
}

variable "vpc_id" {
  description = "ID da VPC onde o RDS será criado"
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR da VPC (para liberar acesso no Security Group)"
  type        = string
}

variable "private_subnet_ids" {
  description = "Lista de IDs das subnets privadas"
  type        = list(string)
}

# Banco
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
  description = "Senha mestre do banco (use uma variável sensível)"
  type        = string
  sensitive   = true
}

# Engine
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

# Instância
variable "instance_class" {
  description = "Classe da instância (ex: db.t4g.small)"
  type        = string
  default     = "db.t4g.small"
}

variable "instance_count" {
  description = "Número de instâncias no cluster (1 para dev, 2 para prod)"
  type        = number
  default     = 1
}

# Escalabilidade Serverless
variable "min_capacity" {
  description = "Capacidade mínima (ACUs) para Serverless v2"
  type        = number
  default     = 0.5
}

variable "max_capacity" {
  description = "Capacidade máxima (ACUs) para Serverless v2"
  type        = number
  default     = 4
}

# Backup
variable "backup_retention_days" {
  description = "Dias de retenção de backup (0 desabilita)"
  type        = number
  default     = 7
}

# Proteção
variable "deletion_protection" {
  description = "Proteger o cluster contra deleção acidental"
  type        = bool
  default     = false
}

variable "apply_immediately" {
  description = "Aplicar mudanças imediatamente (true) ou na janela de manutenção"
  type        = bool
  default     = false
}

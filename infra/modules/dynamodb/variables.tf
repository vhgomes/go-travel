variable "project" {
  description = "Nome do projeto"
  type        = string
}

variable "environment" {
  description = "Ambiente (dev, staging, prod)"
  type        = string
}

variable "deletion_protection" {
  description = "Proteger tabela contra deleção acidental"
  type        = bool
  default     = false
}

variable "ttl_enabled" {
  description = "Habilitar TTL (expiração automática) na tabela principal"
  type        = bool
  default     = true
}

variable "enable_audit" {
  description = "Criar tabela de auditoria"
  type        = bool
  default     = false
}

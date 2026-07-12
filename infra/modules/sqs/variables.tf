variable "project" {
  description = "Nome do projeto"
  type        = string
}

variable "environment" {
  description = "Ambiente (dev, staging, prod)"
  type        = string
}

variable "max_receive_count" {
  description = "Número máximo de tentativas antes de enviar para DLQ"
  type        = number
  default     = 3
}

variable "enable_notifications" {
  description = "Criar fila de notificações (não FIFO)"
  type        = bool
  default     = true
}

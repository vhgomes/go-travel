# Tabela principal de estado do Saga
resource "aws_dynamodb_table" "saga_state" {
  name           = "${var.project}-${var.environment}-saga-state"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "order_id"
  range_key      = "step"

  attribute {
    name = "order_id"
    type = "S"
  }

  attribute {
    name = "step"
    type = "S"
  }

  # GSI para consultar por status (ex: todos os pedidos com erro)
  attribute {
    name = "status"
    type = "S"
  }

  attribute {
    name = "updated_at"
    type = "S"
  }

  global_secondary_index {
    name               = "status-index"
    hash_key           = "status"
    range_key          = "updated_at"
    projection_type    = "ALL"
  }

  # TTL para limpeza automática após 30 dias (opcional)
  ttl {
    attribute_name = "ttl"
    enabled        = var.ttl_enabled
  }

  # Proteção contra deleção acidental (opcional via variável)
  deletion_protection_enabled = var.deletion_protection

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-saga-state"
  }
}

# Tabela de auditoria (opcional - para rastrear mudanças de estado)
resource "aws_dynamodb_table" "saga_audit" {
  count = var.enable_audit ? 1 : 0

  name           = "${var.project}-${var.environment}-saga-audit"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "order_id"
  range_key      = "timestamp"

  attribute {
    name = "order_id"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "S"
  }

  attribute {
    name = "status"
    type = "S"
  }

  global_secondary_index {
    name               = "status-timestamp-index"
    hash_key           = "status"
    range_key          = "timestamp"
    projection_type    = "ALL"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-saga-audit"
  }
}

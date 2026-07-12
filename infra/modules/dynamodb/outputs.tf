output "saga_state_table_name" {
  value = aws_dynamodb_table.saga_state.name
}

output "saga_state_table_arn" {
  value = aws_dynamodb_table.saga_state.arn
}

output "saga_state_table_id" {
  value = aws_dynamodb_table.saga_state.id
}

output "audit_table_name" {
  value = var.enable_audit ? aws_dynamodb_table.saga_audit[0].name : null
}

output "audit_table_arn" {
  value = var.enable_audit ? aws_dynamodb_table.saga_audit[0].arn : null
}

output "cluster_id" {
  value = aws_rds_cluster.this.id
}

output "cluster_endpoint" {
  value = aws_rds_cluster.this.endpoint
}

output "cluster_reader_endpoint" {
  value = aws_rds_cluster.this.reader_endpoint
}

output "cluster_arn" {
  value = aws_rds_cluster.this.arn
}

output "instance_endpoint" {
  value = aws_rds_cluster_instance.this[0].endpoint
}

output "security_group_id" {
  value = aws_security_group.rds.id
}

output "database_name" {
  value = var.database_name
}

output "master_username" {
  value = var.master_username
}

output "master_password" {
  value     = var.master_password
  sensitive = true
}

output "connection_string" {
  value     = "postgres://${var.master_username}:${var.master_password}@${aws_rds_cluster.this.endpoint}/${var.database_name}"
  sensitive = true
}

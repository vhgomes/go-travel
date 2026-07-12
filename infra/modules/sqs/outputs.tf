# URLs das filas principais
output "saga_start_queue_url" {
  value = aws_sqs_queue.saga_start.url
}

output "flight_book_queue_url" {
  value = aws_sqs_queue.flight_book.url
}

output "hotel_book_queue_url" {
  value = aws_sqs_queue.hotel_book.url
}

output "payment_process_queue_url" {
  value = aws_sqs_queue.payment_process.url
}

output "rollback_queue_url" {
  value = aws_sqs_queue.rollback.url
}

output "notifications_queue_url" {
  value = var.enable_notifications ? aws_sqs_queue.notifications[0].url : null
}

# ARNs das filas principais
output "saga_start_queue_arn" {
  value = aws_sqs_queue.saga_start.arn
}

output "flight_book_queue_arn" {
  value = aws_sqs_queue.flight_book.arn
}

output "hotel_book_queue_arn" {
  value = aws_sqs_queue.hotel_book.arn
}

output "payment_process_queue_arn" {
  value = aws_sqs_queue.payment_process.arn
}

output "rollback_queue_arn" {
  value = aws_sqs_queue.rollback.arn
}

# URLs das DLQs
output "saga_start_dlq_url" {
  value = aws_sqs_queue.saga_start_dlq.url
}

output "flight_book_dlq_url" {
  value = aws_sqs_queue.flight_book_dlq.url
}

output "hotel_book_dlq_url" {
  value = aws_sqs_queue.hotel_book_dlq.url
}

output "payment_process_dlq_url" {
  value = aws_sqs_queue.payment_process_dlq.url
}

output "rollback_dlq_url" {
  value = aws_sqs_queue.rollback_dlq.url
}

# Mapas para facilitar iteração no código
output "queue_urls" {
  value = {
    saga_start    = aws_sqs_queue.saga_start.url
    flight_book   = aws_sqs_queue.flight_book.url
    hotel_book    = aws_sqs_queue.hotel_book.url
    payment       = aws_sqs_queue.payment_process.url
    rollback      = aws_sqs_queue.rollback.url
    notifications = var.enable_notifications ? aws_sqs_queue.notifications[0].url : null
  }
}

output "queue_arns" {
  value = {
    saga_start    = aws_sqs_queue.saga_start.arn
    flight_book   = aws_sqs_queue.flight_book.arn
    hotel_book    = aws_sqs_queue.hotel_book.arn
    payment       = aws_sqs_queue.payment_process.arn
    rollback      = aws_sqs_queue.rollback.arn
    notifications = var.enable_notifications ? aws_sqs_queue.notifications[0].arn : null
  }
}

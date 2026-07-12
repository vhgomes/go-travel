output "vpc_id" {
  value = module.networking.vpc_id
}

output "saga_start_queue_url" {
  value = module.sqs.saga_start_queue_url
}

output "flight_book_queue_url" {
  value = module.sqs.flight_book_queue_url
}

output "hotel_book_queue_url" {
  value = module.sqs.hotel_book_queue_url
}

output "payment_process_queue_url" {
  value = module.sqs.payment_process_queue_url
}

output "rollback_queue_url" {
  value = module.sqs.rollback_queue_url
}

output "notifications_queue_url" {
  value = module.sqs.notifications_queue_url
}

output "queue_urls" {
  value = module.sqs.queue_urls
}

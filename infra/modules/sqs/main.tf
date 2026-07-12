# ============================================================
# 1. Filas Principais (FIFO) + DLQs correspondentes
# ============================================================

# 1.1 Saga Start Queue (inicia o fluxo do pedido)
resource "aws_sqs_queue" "saga_start" {
  name = "${var.project}-${var.environment}-saga-start.fifo"
  fifo_queue = true
  content_based_deduplication = true

  visibility_timeout_seconds = 60   # 1 minuto para processar o início do Saga
  message_retention_seconds  = 86400 # 1 dia
  receive_wait_time_seconds  = 10   # Long polling

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.saga_start_dlq.arn
    maxReceiveCount     = var.max_receive_count
  })

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-saga-start.fifo"
  }
}

# DLQ para saga-start
resource "aws_sqs_queue" "saga_start_dlq" {
  name = "${var.project}-${var.environment}-saga-start-dlq.fifo"
  fifo_queue = true
  content_based_deduplication = true

  message_retention_seconds = 1209600 # 14 dias (máximo permitido)

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-saga-start-dlq.fifo"
  }
}

# 1.2 Flight Book Queue
resource "aws_sqs_queue" "flight_book" {
  name = "${var.project}-${var.environment}-flight-book.fifo"
  fifo_queue = true
  content_based_deduplication = true

  visibility_timeout_seconds = 30
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.flight_book_dlq.arn
    maxReceiveCount     = var.max_receive_count
  })

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-flight-book.fifo"
  }
}

# DLQ para flight-book
resource "aws_sqs_queue" "flight_book_dlq" {
  name = "${var.project}-${var.environment}-flight-book-dlq.fifo"
  fifo_queue = true
  content_based_deduplication = true

  message_retention_seconds = 1209600

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-flight-book-dlq.fifo"
  }
}

# 1.3 Hotel Book Queue
resource "aws_sqs_queue" "hotel_book" {
  name = "${var.project}-${var.environment}-hotel-book.fifo"
  fifo_queue = true
  content_based_deduplication = true

  visibility_timeout_seconds = 30
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.hotel_book_dlq.arn
    maxReceiveCount     = var.max_receive_count
  })

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-hotel-book.fifo"
  }
}

# DLQ para hotel-book
resource "aws_sqs_queue" "hotel_book_dlq" {
  name = "${var.project}-${var.environment}-hotel-book-dlq.fifo"
  fifo_queue = true
  content_based_deduplication = true

  message_retention_seconds = 1209600

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-hotel-book-dlq.fifo"
  }
}

# 1.4 Payment Process Queue
resource "aws_sqs_queue" "payment_process" {
  name = "${var.project}-${var.environment}-payment-process.fifo"
  fifo_queue = true
  content_based_deduplication = true

  visibility_timeout_seconds = 45
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.payment_process_dlq.arn
    maxReceiveCount     = var.max_receive_count
  })

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-payment-process.fifo"
  }
}

# DLQ para payment-process
resource "aws_sqs_queue" "payment_process_dlq" {
  name = "${var.project}-${var.environment}-payment-process-dlq.fifo"
  fifo_queue = true
  content_based_deduplication = true

  message_retention_seconds = 1209600

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-payment-process-dlq.fifo"
  }
}

# 1.5 Rollback Queue (para ações de compensação)
resource "aws_sqs_queue" "rollback" {
  name = "${var.project}-${var.environment}-rollback.fifo"
  fifo_queue = true
  content_based_deduplication = true

  visibility_timeout_seconds = 30
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.rollback_dlq.arn
    maxReceiveCount     = var.max_receive_count
  })

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-rollback.fifo"
  }
}

# DLQ para rollback
resource "aws_sqs_queue" "rollback_dlq" {
  name = "${var.project}-${var.environment}-rollback-dlq.fifo"
  fifo_queue = true
  content_based_deduplication = true

  message_retention_seconds = 1209600

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-rollback-dlq.fifo"
  }
}

# ============================================================
# 2. Notifications Queue (não precisa ser FIFO)
# ============================================================
resource "aws_sqs_queue" "notifications" {
  count = var.enable_notifications ? 1 : 0

  name = "${var.project}-${var.environment}-notifications"
  fifo_queue = false

  visibility_timeout_seconds = 30
  message_retention_seconds  = 86400
  receive_wait_time_seconds  = 10

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.notifications_dlq[0].arn
    maxReceiveCount     = var.max_receive_count
  })

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-notifications"
  }
}

# DLQ para notifications
resource "aws_sqs_queue" "notifications_dlq" {
  count = var.enable_notifications ? 1 : 0

  name = "${var.project}-${var.environment}-notifications-dlq"

  message_retention_seconds = 1209600

  tags = {
    Project     = var.project
    Environment = var.environment
    Name        = "${var.project}-${var.environment}-notifications-dlq"
  }
}

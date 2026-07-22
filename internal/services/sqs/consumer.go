package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/vhgomes/go-travel/pkg/logger"
	"go.uber.org/zap"
)

type MessageHandler func(ctx context.Context, message *types.Message) error

type Consumer struct {
	client   *sqs.Client
	queueURL string
}

func NewConsumer(client *sqs.Client, queueURL string) *Consumer {
	return &Consumer{
		client:   client,
		queueURL: queueURL,
	}
}

func (c *Consumer) ReceiveMessages(ctx context.Context, maxMessages int32, waitTimeSeconds int32) ([]types.Message, error) {
	if maxMessages > 10 {
		maxMessages = 10
	}

	input := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.queueURL),
		MaxNumberOfMessages:   maxMessages,
		WaitTimeSeconds:       waitTimeSeconds,
		MessageAttributeNames: []string{"All"},
		AttributeNames:        []types.QueueAttributeName{"All"},
	}

	result, err := c.client.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("receiving messages: %w", err)
	}

	return result.Messages, nil
}

func (c *Consumer) DeleteMessage(ctx context.Context, receiptHandle string) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := c.client.DeleteMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("deleting message: %w", err)
	}

	return nil
}

func (c *Consumer) Poll(ctx context.Context, handler MessageHandler, maxMessages int32, waitTimeSeconds int32) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			messages, err := c.ReceiveMessages(ctx, maxMessages, waitTimeSeconds)
			if err != nil {
				return err
			}

			for _, msg := range messages {
				if err := handler(ctx, &msg); err != nil {
					logger.Error("error handling message", err, zap.String("message_id", aws.ToString(msg.MessageId)))
					continue
				}

				if err := c.DeleteMessage(ctx, *msg.ReceiptHandle); err != nil {
					logger.Error("error deleting message", err, zap.String("message_id", aws.ToString(msg.MessageId)), zap.String("receipt_handle", *msg.ReceiptHandle))
				}
			}
		}
	}
}

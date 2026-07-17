package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type MessageProducer interface {
	SendMessage(ctx context.Context, body string, attributes map[string]string) (string, error)
}

type Producer struct {
	client   *sqs.Client
	queueURL string
}

func NewProducer(client *sqs.Client, queueURL string) *Producer {
	return &Producer{
		client:   client,
		queueURL: queueURL,
	}
}

func (p *Producer) SendMessage(ctx context.Context, body string, attributes map[string]string) (string, error) {
	messageAttributes := make(map[string]types.MessageAttributeValue)

	for key, value := range attributes {
		messageAttributes[key] = types.MessageAttributeValue{
			StringValue: aws.String(value),
			DataType:    aws.String("String"),
		}
	}

	input := &sqs.SendMessageInput{
		QueueUrl:          aws.String(p.queueURL),
		MessageBody:       aws.String(body),
		MessageAttributes: messageAttributes,
	}

	result, err := p.client.SendMessage(ctx, input)
	if err != nil {
		return "", fmt.Errorf("sending message: %w", err)
	}

	return *result.MessageId, nil
}

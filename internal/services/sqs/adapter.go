package sqs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type QueueManager struct {
	client *sqs.Client
}

func NewQueueManager(client *sqs.Client) *QueueManager {
	return &QueueManager{client: client}
}

func (m *QueueManager) CreateStandardQueue(ctx context.Context, name string, visibilityTimeout int, messageRetention int) (string, error) {
	input := &sqs.CreateQueueInput{
		QueueName: aws.String(name),
		Attributes: map[string]string{
			"VisibilityTimeout":             strconv.Itoa(visibilityTimeout),
			"MessageRetentionPeriod":        strconv.Itoa(messageRetention),
			"ReceiveMessageWaitTimeSeconds": "20",
		},
	}

	result, err := m.client.CreateQueue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("creating standard queue: %w", err)
	}

	return *result.QueueUrl, nil
}

func (m *QueueManager) CreateFIFOQueue(ctx context.Context, name string, contentBasedDedup bool) (string, error) {
	fifoName := name + ".fifo"

	attributes := map[string]string{
		"FifoQueue":                     "true",
		"VisibilityTimeout":             "30",
		"ReceiveMessageWaitTimeSeconds": "20",
	}

	if contentBasedDedup {
		attributes["ContentBasedDeduplication"] = "true"
	}

	input := &sqs.CreateQueueInput{
		QueueName:  aws.String(fifoName),
		Attributes: attributes,
	}

	result, err := m.client.CreateQueue(ctx, input)
	if err != nil {
		return "", fmt.Errorf("creating FIFO queue: %w", err)
	}

	return *result.QueueUrl, nil
}

func (m *QueueManager) GetQueueURL(ctx context.Context, name string) (string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	}

	result, err := m.client.GetQueueUrl(ctx, input)
	if err != nil {
		return "", fmt.Errorf("getting queue URL: %w", err)
	}

	return *result.QueueUrl, nil
}

func (m *QueueManager) ConfigureDeadLetterQueue(ctx context.Context, mainQueueURL string, dlqARN string, maxReceiveCount int) error {
	redrivePolicy := fmt.Sprintf(
		`{"deadLetterTargetArn":"%s","maxReceiveCount":"%d"}`,
		dlqARN,
		maxReceiveCount,
	)

	input := &sqs.SetQueueAttributesInput{
		QueueUrl: aws.String(mainQueueURL),
		Attributes: map[string]string{
			"RedrivePolicy": redrivePolicy,
		},
	}

	_, err := m.client.SetQueueAttributes(ctx, input)
	if err != nil {
		return fmt.Errorf("configuring dead letter queue: %w", err)
	}

	return nil
}

func (m *QueueManager) GetQueueARN(ctx context.Context, queueURL string) (string, error) {
	input := &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
	}

	result, err := m.client.GetQueueAttributes(ctx, input)
	if err != nil {
		return "", fmt.Errorf("getting queue ARN: %w", err)
	}

	return result.Attributes["QueueArn"], nil
}

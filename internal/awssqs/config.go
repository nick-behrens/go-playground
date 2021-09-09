package awssqs

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Config struct {
	QueueName string
	Timeout   int32
}

func (c *Config) GetReceiveMessageInput() (*sqs.Client, *sqs.ReceiveMessageInput, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))

	if err != nil {
		log.Fatal().Err(err).Msg("unable to load SDK config")
		return nil, nil, err
	}

	svc := sqs.NewFromConfig(cfg)

	urlResult, err := svc.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: &c.QueueName,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("unable to get queue url")
		return nil, nil, err
	}

	queueUrl := urlResult.QueueUrl

	sqsReceiveConfig := &sqs.ReceiveMessageInput{
		QueueUrl:            queueUrl,
		MaxNumberOfMessages: int32(10),
		VisibilityTimeout:   c.Timeout,
	}

	return svc, sqsReceiveConfig, nil
}

func (c *Config) SendMessage(client *sqs.Client, body *string) (*sqs.SendMessageOutput, error) {
	urlResult, err := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: &c.QueueName,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("unable to get queue url")
		return nil, err
	}

	queueUrl := urlResult.QueueUrl

	sqsSendMessageInput := &sqs.SendMessageInput{
		MessageBody: body,
		QueueUrl:    queueUrl,
	}

	sendMessageOutput, err := client.SendMessage(context.TODO(), sqsSendMessageInput)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to send message to queue")
		return nil, err
	}

	return sendMessageOutput, nil
}

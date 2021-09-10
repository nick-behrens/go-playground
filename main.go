package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/snapdocs/go-common/database"

	"internal/awssqs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	internalDatabase "github.com/snapdocs/go-playground/pkg/database"
	"github.com/snapdocs/go-playground/pkg/database/queries"
	"github.com/snapdocs/go-playground/pkg/models"
)

type MessageData struct {
	Data string `json:"data"`
}

type EmailRequestsPGPool struct {
	Id             int
	Payload        string
	Closing_id     pgtype.TextArray
	To             string
	Status         string
	Status_message string
	Email_provider string
	Message_uuid   string
	Created_at     string
	Updated_at     string
	Sent_at        string
	Cc             string
	Bcc            string
	Reply_to       string
	Company        string
	Tags           string
	Rumi_version   string
	Expires_at     string
	Attemps        string
	Source         string
	Company_id     string
}

type EmailRequests struct {
	Id             int
	Payload        string
	Closing_id     string
	To             string
	Status         string
	Status_message string
	Email_provider string
	Message_uuid   string
	Created_at     string
	Updated_at     string
	Sent_at        string
	Cc             string
	Bcc            string
	Reply_to       string
	Company        string
	Tags           string
	Rumi_version   string
	Expires_at     string
	Attemps        string
	Source         string
	Company_id     string
}

func main() {
	// cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("unable to load SDK config")
	// }

	// svc := sqs.NewFromConfig(cfg)

	// createQueues(svc)

	// sendAndReceiveSQS()

	// connectAndRunQueryWithGorm()

	// connectAndRunQueryWithGoCommon()

	// getDatabaseAndRunQueries()

	// testDatadogAndAirbrakes()

	// simpleChannelExample()

	// channelsToPassDataBetweenRoutines()

	// channelsToPassDataAndEndFunction()

	multipleChannelProducers()
}

func getDatabaseAndRunQueries() {
	adapter, err := internalDatabase.NewDatabase(internalDatabase.Config{
		Host:     "localhost",
		User:     "snapp",
		Password: "snapp",
		DBName:   "go_playground",
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create database")
	}

	queries := queries.New(adapter)

	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create uuid")
	}

	createdWebhookEvent, err := queries.CreateWebhookEvent(&models.WebhookEvent{
		WebhookEvent:           "event",
		EmailProvider:          "provider",
		To:                     "to",
		EmailProviderMessageId: "provider_id",
		Reason:                 "reason",
		Event:                  "event",
		SdMessageId:            uuid,
		Timestamp:              1234,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create webhook event.")
		return
	}

	log.Info().Msg("Created webhookEvent ID is: " + createdWebhookEvent.ID.String())

	webhookByID, _ := queries.GetWebhookByID(createdWebhookEvent.ID)

	log.Info().Msg("Result of getting the webhookEvent by ID: " + webhookByID.ID.String())

	webhookEvent, _ := queries.GetFirstWebhookEvent()

	log.Info().Msg("Result of getting the first webhookEvent: " + webhookEvent.ID.String())
}

func connectAndRunQueryWithGorm() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		"localhost",
		"snapp",
		"snapp",
		"go_playground",
		"5432",
		"disable",
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize gorm connection")
	}

	var test MessageData

	db.Table("testing").Create(&MessageData{
		Data: "I am data from Go!",
	})

	db.Table("testing").First(&test)

	log.Info().Msg("Result of show tables: " + test.Data)
}

func connectAndRunQueryWithGoCommon() {
	dbConfig := &database.DbConfig{
		Adapter:  "postgresql",
		Encoding: "unicode",
		Database: "rumi_development",
		Username: "snapp",
		Password: "snapp",
		Pool:     5,
		Port:     5432,
		Host:     "localhost",
		Timeout:  5000,
	}

	err := dbConfig.InitializeDefaultDBPool()

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to initialize database connection pool")
	}

	connection, err := database.Pool.Acquire(context.TODO())

	if err != nil {
		log.Fatal().Err(err).Msg("Error acquiring a connection")
	}

	row := connection.QueryRow(context.TODO(), "SELECT * FROM email_requests limit 1")

	var i EmailRequestsPGPool

	err = row.Scan(
		&i.Id,
		&i.Payload,
		&i.Closing_id,
		&i.To,
		&i.Status,
		&i.Status_message,
		&i.Email_provider,
		&i.Message_uuid,
		&i.Created_at,
		&i.Updated_at,
		&i.Sent_at,
		&i.Cc,
		&i.Bcc,
		&i.Reply_to,
		&i.Company,
		&i.Tags,
		&i.Rumi_version,
		&i.Expires_at,
		&i.Attemps,
		&i.Source,
		&i.Company_id,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Row scan failed.")
	}

	log.Info().Msg("Result of show tables: " + i.Message_uuid)
}

func sendAndReceiveSQS() {
	awscfg := &awssqs.Config{
		QueueName: "development-nickbehrens-rumi",
		Timeout:   60,
	}

	svc, sqsReceiveConfig, err := awscfg.GetReceiveMessageInput()

	if err != nil {
		log.Fatal().Err(err).Msg("unable to load message input")
		return
	}

	testMessage := map[string]string{"data": "I'm a message!"}
	testMessageJSON, _ := json.Marshal(testMessage)
	testString := string(testMessageJSON)

	sendMessageOutput, err := awscfg.SendMessage(svc, &testString)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to send message")
		return
	}

	log.Info().Msg("Message Id of sent message: " + *sendMessageOutput.MessageId)

	receiveMessageOutput, err := svc.ReceiveMessage(context.TODO(), sqsReceiveConfig)

	if err != nil {
		log.Fatal().Err(err).Msg("unable to receive message")
		return
	}

	res := MessageData{}

	for _, message := range receiveMessageOutput.Messages {
		err = json.Unmarshal([]byte(*message.Body), &res)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not unmarsall data to type")
		} else {
			log.Info().Msg("Message Body was: " + res.Data)
		}
	}
}

func createQueues(sqsClient *sqs.Client) {
	var sb strings.Builder

	queue := "development-nickbehrens-rumi"

	// Make DLQ
	dlqAttributes := make(map[string]string)
	dlqAttributes["MessageRetentionPeriod"] = strconv.FormatInt(60*60*24*14, 10)

	sb.WriteString(queue)
	sb.WriteString("-dlq")

	dlqName := sb.String()

	createdDLQ, err := sqsClient.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName:  &dlqName,
		Attributes: dlqAttributes,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create dlq")
	}

	log.Info().Str("Queue url for DLQ: ", *createdDLQ.QueueUrl).Msg("DLQ url is in this message")

	dlqAttributesAWS, err := sqsClient.GetQueueAttributes(context.TODO(), &sqs.GetQueueAttributesInput{
		QueueUrl: createdDLQ.QueueUrl,
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameQueueArn,
		},
	})

	if err != nil {
		log.Fatal().Err(err).Msg("failed to get dlq attributes")
	}

	// Make main queue.
	mainQueueAttributes := make(map[string]string)
	redrivePolicy := map[string]string{"maxReceiveCount": "3", "deadLetterTargetArn": dlqAttributesAWS.Attributes["QueueArn"]}
	redrivePolicyJSON, _ := json.Marshal(redrivePolicy)
	mainQueueAttributes["RedrivePolicy"] = string(redrivePolicyJSON)
	mainQueueAttributes["VisibilityTimeout"] = "65"
	mainQueueAttributes["MessageRetentionPeriod"] = strconv.FormatInt(60*60*24*4, 10)
	mainQueueAttributes["DelaySeconds"] = "0"
	mainQueueAttributes["MaximumMessageSize"] = strconv.FormatInt(256*1024, 10)
	mainQueueAttributes["ReceiveMessageWaitTimeSeconds"] = "0"

	// If the queue is the same as an existing one, the method will return and not have any issues.
	createdQueue, err := sqsClient.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName:  &queue,
		Attributes: mainQueueAttributes,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create queue")
	}

	log.Info().Str("Queue url for Main Queue: ", *createdQueue.QueueUrl).Msg("DLQ url is in this message")
}

package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"log"
	"os"
)

type Tag struct {
	Key   string
	Value string
}

type Detail struct {
	SecretId string `json:"secretId"`
	Tags     []Tag  `json:"tags"`
}

func LambdaHandler(ctx context.Context, secretName string) {
	sourceName := os.Getenv("SOURCE")
	eventBusName := os.Getenv("EventBus")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		log.Fatal(err)
	}

	events := cloudwatchevents.New(sess)
	entries := []*cloudwatchevents.PutEventsRequestEntry{}

	tags := getSecretTags(secretName, sess)

	detail := Detail{
		SecretId: secretName,
		Tags:     tags,
	}
	detailB, _ := json.Marshal(detail)
	detailS := string(detailB)
	detailType := "SecretId, Tags"

	entry := cloudwatchevents.PutEventsRequestEntry{
		Detail:       &detailS,
		DetailType:   &detailType,
		Source:       &sourceName,
		EventBusName: &eventBusName,
	}

	entries = append(entries, &entry)
	input := cloudwatchevents.PutEventsInput{Entries: entries}
	res, err := events.PutEvents(&input)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.String())
}

func getSecretTags(secretId string, s *session.Session) []Tag {
	tags := []Tag{}
	secrets := secretsmanager.New(s)
	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretId),
	}

	output, err := secrets.DescribeSecret(input)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range output.Tags {
		tag := Tag{
			Key:   *t.Key,
			Value: *t.Value,
		}
		tags = append(tags, tag)
	}

	return tags
}

func main() {
	lambda.Start(LambdaHandler)
}

package utils

import (
	"../config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

var (
	AwsSession  *session.Session
	AwsKinesis  *kinesis.Kinesis
	AwsFirehose *firehose.Firehose
)

func CreateAwsSession() {
	creds := credentials.NewStaticCredentials(config.Config.AwsKey, config.Config.AwsSecret, "")
	_, err := creds.Get()

	if err != nil {
		panic(err)
	}

	AwsSession, err = session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: creds,
	})

	AwsKinesis = kinesis.New(AwsSession)
	AwsFirehose = firehose.New(AwsSession)

	if err != nil {
		panic(err)
	}
}

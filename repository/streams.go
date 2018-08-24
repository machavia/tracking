package repository

import (
	"../config"
	"../utils"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/firehose"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"strconv"
	"sync"
)

const BufferSize = 490

var mu sync.Mutex
var hitBuffer = make([]*firehose.Record, 0, BufferSize)

//sending the hit data as json string in Kinesis Stream for real time prediction
//the data are splitted into shards by using the visitor id partition key. Which means, a same visitor will always be in the same shard.
func SendToStream(visitorId uint64, hit []byte) {

	_, err := utils.AwsKinesis.PutRecord(&kinesis.PutRecordInput{
		Data:         hit,
		StreamName:   aws.String(config.Config.KinesisStreamName),
		PartitionKey: aws.String(strconv.FormatUint(visitorId, 10)),
	})
	if err != nil {
		fmt.Println("kinesis error")
		panic(err)
	}

}

//sending the hit data as json string in Kinesis Stream for real time prediction
//this function act as a buffer. The data are sent to Kinsesis Firehose only when the buffer is full.
//We do that in order to minimize the number of request across the network, and the performance of the application
func SendToFirehose(hit []byte) {
	mu.Lock()
	defer mu.Unlock()
	hitBuffer = append(
		hitBuffer,
		&firehose.Record{Data: append(hit, '\n')},
	)

	//if the buffer if not full we continue to add hits
	if len(hitBuffer) < BufferSize {
		return
	}

	//else we want to send all the hits we buffered to firehose
	_, err := utils.AwsFirehose.PutRecordBatch(&firehose.PutRecordBatchInput{
		DeliveryStreamName: aws.String(config.Config.FirehoseStreamName),
		Records:            hitBuffer,
	},
	)

	if err != nil {
		fmt.Println("Buffer Size", len(hitBuffer))
		panic(err)
	}

	fmt.Println("-> sending to firehose")

	//emptying the slice for the next batch
	hitBuffer = make([]*firehose.Record, 0, BufferSize)
}

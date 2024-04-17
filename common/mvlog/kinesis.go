package mvlog

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

type KinesisConfig struct {
	Region     string
	AccessKey  string
	SecretKey  string
	StreamName string
}

var kinesisCfg KinesisConfig
var kinesisClient *kinesis.Client
var shardId string

func InitKinesis(ctx context.Context, cfg KinesisConfig) {
	kinesisCfg = cfg
	log.Println(cfg)
	if cfg.AccessKey != "" {
		awsConfig, _ := config.LoadDefaultConfig(
			ctx,
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
			config.WithRegion(cfg.Region),
		)

		kinesisClient = kinesis.NewFromConfig(awsConfig)

	} else {
		awsConfig, _ := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))

		kinesisClient = kinesis.NewFromConfig(awsConfig)
	}

	stream, err := kinesisClient.DescribeStream(
		ctx,
		&kinesis.DescribeStreamInput{StreamName: &cfg.StreamName})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("kiness connect success ", *stream.StreamDescription.StreamName)
}

var kinesisRecords []types.PutRecordsRequestEntry

func BeginKinesis() {
	kinesisRecords = make([]types.PutRecordsRequestEntry, 0)
}

func PutKinesisMulti(pk string, data []byte) {
	if kinesisClient == nil {
		return
	}

	kinesisRecords = append(kinesisRecords, types.PutRecordsRequestEntry{
		Data:         data,
		PartitionKey: aws.String(pk),
	})
}

func PutKinesis(ctx context.Context, pk string, data []byte) error {
	out, err := kinesisClient.PutRecord(ctx,
		&kinesis.PutRecordInput{
			Data:         data,
			StreamName:   &kinesisCfg.StreamName,
			PartitionKey: aws.String(pk),
		})
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Println(*out.ShardId)
	shardId = *out.ShardId
	return err
}

func EndKinesis(ctx context.Context) error {
	if len(kinesisRecords) == 0 {
		return nil
	}
	_, err := kinesisClient.PutRecords(ctx,
		&kinesis.PutRecordsInput{
			Records:    kinesisRecords,
			StreamName: &kinesisCfg.StreamName,
		})
	if err != nil {
		log.Println(err)
		return err
	}
	// for _, v := range out.Records {
	// 	log.Println(*v.ShardId)
	// 	//shardId = out.ShardId
	// }
	return err
}

func PopKinesis(ctx context.Context) {
	// retrieve iterator
	iteratorOutput, err := kinesisClient.GetShardIterator(ctx,
		&kinesis.GetShardIteratorInput{
			// Shard Id is provided when making put record(s) request.
			ShardId:           aws.String("shardId-000000000002"),
			ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
			//aws.String("TRIM_HORIZON"),
			// ShardIteratorType: aws.String("AT_SEQUENCE_NUMBER"),
			// ShardIteratorType: aws.String("LATEST"),
			StreamName: &kinesisCfg.StreamName,
		})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(iteratorOutput)

	// get records use shard iterator for making request
	records, err := kinesisClient.GetRecords(ctx,
		&kinesis.GetRecordsInput{
			ShardIterator: iteratorOutput.ShardIterator,
		})
	if err != nil {
		log.Println(err)
		return
	}
	for _, record := range records.Records {
		log.Println("record: ", record)
		log.Println("record-Data: ", string(record.Data))
	}

	log.Println(*records.NextShardIterator)
}

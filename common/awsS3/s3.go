package awss3

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Config ...
type Config struct {
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

var client *s3.Client

func Init(cfg Config) {
	if cfg.AccessKey != "" {
		awsConfig, _ := config.LoadDefaultConfig(
			context.Background(),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
			config.WithRegion(cfg.Region),
		)
		client = s3.NewFromConfig(awsConfig)
	} else {
		awsConfig, _ := config.LoadDefaultConfig(
			context.Background(),
			config.WithRegion(cfg.Region),
		)

		client = s3.NewFromConfig(awsConfig)
	}

	log.Println("s3 connect success")
}

func ListObjects(ctx context.Context, bucket, path string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &path,
	}

	result, err := client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}

	list := []string{}
	for _, obj := range result.Contents {
		list = append(list, *obj.Key)
	}
	return list, nil
}

func GetObject(ctx context.Context, bucket, name string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &name,
	}
	result, err := client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}

	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}

func PutObject(ctx context.Context, bucket, target string, body []byte) error {
	reader := bytes.NewReader(body)
	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &target,
		Body:   reader,
	}
	_, err := client.PutObject(ctx, input)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DelObject(ctx context.Context, bucket, target string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &target,
	}

	_, err := client.DeleteObject(ctx, input)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DelObjects(ctx context.Context, bucket, target string) error {
	list, err := ListObjects(ctx, bucket, target)
	if err != nil {
		log.Println(err)
		return err
	}
	list = append(list, target)

	for _, v := range list {
		log.Println("del --- ", v)
		input := &s3.DeleteObjectInput{
			Bucket: &bucket,
			Key:    &v,
		}
		_, err := client.DeleteObject(ctx, input)
		if err != nil {
			log.Println(err)
			return err
		}

	}
	return nil
}

package awsssm

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/chipong/template/core"
	"gopkg.in/yaml.v3"
)

// Config ...
type Config struct {
	Region    string
	AccessKey string
	SecretKey string
}

var client *ssm.Client

func Init(cfg Config) {
	if cfg.AccessKey != "" {
		awsConfig, _ := config.LoadDefaultConfig(
			context.Background(),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
			config.WithRegion(cfg.Region),
		)
		client = ssm.NewFromConfig(awsConfig)
	} else {
		awsConfig, _ := config.LoadDefaultConfig(
			context.Background(),
			config.WithRegion(cfg.Region),
		)

		client = ssm.NewFromConfig(awsConfig)
	}

	log.Println("ssm connect success")
}

func GetValue(ctx context.Context, name string) (string, error) {
	input := &ssm.GetParameterInput{
		Name: &name,
	}

	parameter, err := client.GetParameter(ctx, input)
	if err != nil {
		return "", err
	}

	value := *parameter.Parameter.Value
	return value, nil
}

func LoadParamaterStore(ctx context.Context, conf interface{}, name string) string {
	v, err := GetValue(ctx, name)
	if err != nil {
		log.Println(err)
		return ""
	}
	err = yaml.Unmarshal([]byte(v), conf)
	if err != nil {
		log.Println(err)
		return ""
	}
	checkSum := core.FileCheckSumOverload(v)
	return checkSum
}

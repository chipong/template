package dynamodb

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// 1 month
var ttl = int64(3600 * 24 * 30)

var dynamodbconn *dynamodb.Client

// Config ...
type Config struct {
	Region    string
	AccessKey string
	SecretKey string
	EndPoint  string
	TTL       int64
}

var tableName = "game-oz"
var dynamodbCfg Config

// InitDynamoDB ...
func InitDynamoDB(cfg Config) {
	dynamodbCfg = cfg
	if cfg.AccessKey != "" {
		awsConfig, _ := config.LoadDefaultConfig(
			context.Background(),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
			config.WithRegion(cfg.Region),
		)

		dynamodbconn = dynamodb.NewFromConfig(awsConfig)

	} else {
		awsConfig, _ := config.LoadDefaultConfig(
			context.Background(),
			config.WithRegion(cfg.Region),
		)

		dynamodbconn = dynamodb.NewFromConfig(awsConfig)
	}

	go ping()
}

// CheckDynamoDBConn ...
func CheckDynamoDBConn(cfg Config) {
	if cfg.Region != dynamodbCfg.Region ||
		cfg.AccessKey != dynamodbCfg.AccessKey ||
		cfg.SecretKey != dynamodbCfg.SecretKey ||
		cfg.EndPoint != dynamodbCfg.EndPoint {
		InitDynamoDB(cfg)
	}
}

func ping() {
	//input := &dynamodb.ListTablesInput{}

	ctx := context.Background()
	var cancelFn func()
	ctx, cancelFn = context.WithTimeout(ctx, time.Second*5)
	defer cancelFn()

	//List the tabkes in this account
	// result, err := dynamodbconn.ListTables(ctx, input)
	// if err != nil {
	// 	log.Println("dynamodb connection fail")
	// 	return
	// }
	//log.Println("dynamodb tables: ", result)
	log.Println("dynamodb connection success")

	// createdTable := false
	// for _, v := range result.TableNames {
	// 	if *v == tableName {
	// 		createdTable = true
	// 	}
	// }

	// if createdTable == false {
	// 	CreateTable()
	// }
}

// CreateTable ...
// func CreateTable() error {
// 	input := &dynamodb.CreateTableInput{
// 		AttributeDefinitions: []*dynamodb.AttributeDefinition{
// 			{
// 				AttributeName: aws.String("uid"),
// 				AttributeType: aws.String("S"),
// 			},
// 		},
// 		KeySchema: []*dynamodb.KeySchemaElement{
// 			{
// 				AttributeName: aws.String("uid"),
// 				KeyType:       aws.String("HASH"),
// 			},
// 		},
// 		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
// 			ReadCapacityUnits:  aws.Int64(5),
// 			WriteCapacityUnits: aws.Int64(5),
// 		},
// 		TableName: aws.String(tableName),
// 	}

// 	_, err := dynamodbconn.CreateTable(input)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	return nil
// }

func BeginTransaction() *dynamodb.TransactWriteItemsInput {
	return &dynamodb.TransactWriteItemsInput{
		TransactItems: make([]types.TransactWriteItem, 0),
	}
}

func AddTransaction(tx *dynamodb.TransactWriteItemsInput, item types.TransactWriteItem) {
	tx.TransactItems = append(tx.TransactItems, item)
}

func EndTransaction(ctx context.Context, tx *dynamodb.TransactWriteItemsInput) error {
	log.Println("tx count: ", len(tx.TransactItems))
	_, err := dynamodbconn.TransactWriteItems(ctx, tx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func MarshalList(in interface{}) ([]types.AttributeValue, error) {
	av, err := attributevalue.NewEncoder(func(eo *attributevalue.EncoderOptions) {
		eo.TagKey = `json`
	}).Encode(in)
	log.Println(av)
	asMap, ok := av.(*types.AttributeValueMemberL)
	if err != nil || av == nil || !ok {
		return []types.AttributeValue{}, err
	}

	return asMap.Value, nil
}

func MarshalMap(in interface{}) (map[string]types.AttributeValue, error) {
	av, err := attributevalue.NewEncoder(func(eo *attributevalue.EncoderOptions) {
		eo.TagKey = `json`
	}).Encode(in)

	asMap, ok := av.(*types.AttributeValueMemberM)
	if err != nil || av == nil || !ok {
		return map[string]types.AttributeValue{}, err
	}

	return asMap.Value, nil
}

func Marshal(in interface{}) (types.AttributeValue, error) {
	return attributevalue.NewEncoder(func(eo *attributevalue.EncoderOptions) {
		eo.TagKey = `json`
	}).Encode(in)
}

func UnmarshalMap(m map[string]types.AttributeValue, out interface{}) error {
	return attributevalue.NewDecoder(func(do *attributevalue.DecoderOptions) {
		do.TagKey = `json`
	}).Decode(&types.AttributeValueMemberM{Value: m}, out)
}

func BeginPQLTransaction() *dynamodb.ExecuteTransactionInput {
	return &dynamodb.ExecuteTransactionInput{
		TransactStatements:     make([]types.ParameterizedStatement, 0),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}
}

func AddPQLTransaction(tx *dynamodb.ExecuteTransactionInput, item types.ParameterizedStatement) {
	tx.TransactStatements = append(tx.TransactStatements, item)
}

func MaxExecPQLTransaction(ctx context.Context, tx *dynamodb.ExecuteTransactionInput) error {
	if len(tx.TransactStatements) == 0 {
		return nil
	}

	for i := 0; i < len(tx.TransactStatements); i++ {
		input := dynamodb.ExecuteStatementInput{
			Statement:  tx.TransactStatements[i].Statement,
			Parameters: tx.TransactStatements[i].Parameters,
		}

		_, err := dynamodbconn.ExecuteStatement(ctx, &input)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	tx.TransactStatements = make([]types.ParameterizedStatement, 0)
	return nil
}

func EndPQLTransaction(ctx context.Context, tx *dynamodb.ExecuteTransactionInput) error {
	log.Println("tx count: ", len(tx.TransactStatements))
	if len(tx.TransactStatements) == 0 {
		return nil
	}

	if len(tx.TransactStatements) == 1 {
		input := dynamodb.ExecuteStatementInput{
			Statement:  tx.TransactStatements[0].Statement,
			Parameters: tx.TransactStatements[0].Parameters,
		}
		_, err := dynamodbconn.ExecuteStatement(ctx, &input)
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	}

	out, err := dynamodbconn.ExecuteTransaction(ctx, tx)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range out.ConsumedCapacity {
		log.Println("tx capacity(", *v.TableName, ", ", *v.WriteCapacityUnits, ")")
	}
	return nil
}

func IntToStr[T int | int32 | int64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

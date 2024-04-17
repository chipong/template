package dynamodb

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/chipong/template/common/proto"
)

func GetTemplates(ctx context.Context, pk string) ([]*oz.OZTemplate, error) {
	query := `SELECT * FROM "game-oz" WHERE "PK" = ? AND BEGINS_WITH("SK", 'TEMPLATE#')`
	param := dynamodb.ExecuteStatementInput{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: pk},
		},
	}
	result, err := dynamodbconn.ExecuteStatement(ctx, &param)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	templates := []*oz.OZTemplate{}
	for _, v := range result.Items {
		template := &oz.OZTemplate{}
		sk := v["SK"].(*types.AttributeValueMemberS).Value
		UnmarshalMap(v, template)
		template.Id = sk[len("TEMPLATE#"):]
		templates = append(templates, template)
	}
	return templates, nil
}

func PutTemplateTx(pk string, template *oz.OZTemplate) types.ParameterizedStatement {
	query := `INSERT INTO "game-oz"
	VALUE {
		'PK':?, 
		'SK':?,
		'count':?,
		'update_at':?,
		'create_at':?
	}`
	at := time.Now().Unix()

	return types.ParameterizedStatement{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: pk},
			&types.AttributeValueMemberS{Value: "TEMPLATE#" + template.Id},
			&types.AttributeValueMemberN{Value: IntToStr(template.Count)},
			&types.AttributeValueMemberN{Value: IntToStr(at)},
			&types.AttributeValueMemberN{Value: IntToStr(at)},
		},
	}
}

func UpdateTemplateCountTx(pk, sk string, count int64) types.ParameterizedStatement {
	query := `UPDATE "game-oz"
		SET "count" = ?,
			"update_at" = ?
	WHERE "PK" = ? AND "SK" = ?`
	at := time.Now().Unix()
	return types.ParameterizedStatement{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberN{Value: IntToStr(count)},
			&types.AttributeValueMemberN{Value: IntToStr(at)},
			&types.AttributeValueMemberS{Value: pk},
			&types.AttributeValueMemberS{Value: sk},
		},
	}
}
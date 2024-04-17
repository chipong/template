package dynamodb

import (
	"context"
	"fmt"
	"log"
	// "strconv"
	// "strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	// "github.com/chipong/template/core"
	// "github.com/chipong/template/common/proto"
)

func BeginPQLBatch() *dynamodb.BatchExecuteStatementInput {
	return &dynamodb.BatchExecuteStatementInput{
		Statements: make([]types.BatchStatementRequest, 0),
	}
}

func AddPQLBatch(batch *dynamodb.BatchExecuteStatementInput, item types.BatchStatementRequest) {
	batch.Statements = append(batch.Statements, item)
}
func EndPQLBatch(ctx context.Context, batch *dynamodb.BatchExecuteStatementInput) ([](interface{}), error) {
	result, err := dynamodbconn.BatchExecuteStatement(ctx, batch)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	items := [](interface{}){}

	for _, i := range result.Responses {
		item := map[string]interface{}{}
		err = UnmarshalMap(i.Item, &item)
		if err != nil {
			log.Println(err)
			continue
		}
		items = append(items, &item)
	}

	return items, nil
}

func GetBatchTable(ctx context.Context, pk string, sk map[string]string) ([]interface{}, error) {
	query := fmt.Sprintf(`SELECT * FROM "game-oz" WHERE "PK" = '%s'`, pk)
	condition := ""
	for k, v := range sk {
		if v == "=" {
			condition += fmt.Sprintf(` "SK" = '%s' OR`, k)
		} else {
			condition += fmt.Sprintf(` BEGINS_WITH("SK", '%s') OR`, k)
		}
	}
	condition = condition[:len(condition)-2]
	query += fmt.Sprintf(` AND (%s)`, condition)

	//log.Println(query)

	param := dynamodb.ExecuteStatementInput{
		Statement:              aws.String(query),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}
	result, err := dynamodbconn.ExecuteStatement(context.Background(), &param)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 아직 검색할 항목이 남아있을 때 nextToken이 return, param에 추가하여 추가 검색
	executeNextToken := result.NextToken
	var nextTokenResult *dynamodb.ExecuteStatementOutput
	for {
		if executeNextToken != nil {
			param.NextToken = executeNextToken

			nextTokenResult, err = dynamodbconn.ExecuteStatement(context.Background(), &param)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			result.Items = append(result.Items, nextTokenResult.Items...)
			executeNextToken = nextTokenResult.NextToken
		} else {
			break
		}
	}

	if result.ConsumedCapacity != nil {
		consumed := *result.ConsumedCapacity
		log.Println("read capacity: ", *consumed.TableName, " -> ", *consumed.CapacityUnits)
	}

	items := [](interface{}){}

	for _, i := range result.Items {
		log.Println(i)
		// sk := i["SK"].(*types.AttributeValueMemberS).Value
		// if strings.Contains(sk, "ACCOUNT#") {
		// 	item := oz.OZAccount{}
		// 	err = UnmarshalMap(i, &item)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "KINGDOM#") {
		// 	item := oz.OZKingdom{}
		// 	err = UnmarshalMap(i, &item)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "ALLIANCE#") {
		// 	item := oz.OZAlliance{}
		// 	err = UnmarshalMap(i, &item)
		// 	item.Nation = sk[len("ALLIANCE#"):]
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "AUTH#") {
		// 	item := Auth{}
		// 	err = UnmarshalMap(i, &item)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "CHAR#") {
		// 	item := oz.OZChar{}
		// 	err = UnmarshalMap(i, &item)
		// 	item.CharId = sk[len("CHAR#"):]
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "SEPHIRA#") {
		// 	item := oz.OZSephira{}
		// 	err = UnmarshalMap(i, &item)
		// 	item.SephiraId = sk[len("SEPHIRA#"):]
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "EQUIP#PAGE#") {
		// 	item := oz.OZEquipInfoes{}
		// 	intSk, _ := strconv.Atoi(sk[11:])
		// 	invenPage := int32(intSk)
		// 	err = UnmarshalMap(i, &item)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	item.InvenIndex = invenPage
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "PARTY#") {
		// 	item := oz.OZParty{}
		// 	err = UnmarshalMap(i, &item)
		// 	partyNo, _ := core.ParseInt(sk[len("PARTY#"):])
		// 	item.PartyNo = int32(partyNo)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "POINT#") {
		// 	item := oz.OZPointAssets{}
		// 	err = UnmarshalMap(i, &item)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "STACK#PAGE#") {
		// 	item := oz.OZStackInfoes{}
		// 	intSk, _ := strconv.Atoi(sk[11:])
		// 	invenPage := int32(intSk)
		// 	err = UnmarshalMap(i, &item)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	item.InvenIndex = invenPage
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "STAGE#") {
		// 	item := oz.OZStage{}
		// 	err = UnmarshalMap(i, &item)
		// 	item.StageId = sk[len("STAGE#"):]
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "QUEST#") {
		// 	item := oz.OZQuestInfo{}
		// 	err = UnmarshalMap(i, &item)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	questType := sk[len("QUEST#"):]
		// 	item.QuestType = questType
		// 	items = append(items, &item)
		// } else if strings.Contains(sk, "BUILDING#PAGE#") {
		// 	item := oz.OZBuildingInfo{}
		// 	intSk, _ := strconv.Atoi(sk[len("BUILDING#PAGE#"):])
		// 	invenPage := int32(intSk)
		// 	err = UnmarshalMap(i, &item)
		// 	if err != nil {
		// 		log.Println(err)
		// 		continue
		// 	}
		// 	item.InvenIndex = invenPage
		// 	items = append(items, &item)
		// }
	}

	return items, nil
}

func DeleteTx(pk, sk string) types.ParameterizedStatement {
	query := `DELETE FROM "game-oz"
		WHERE "PK" = ? AND "SK" = ?`

	return types.ParameterizedStatement{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: pk},
			&types.AttributeValueMemberS{Value: sk},
		},
	}
}

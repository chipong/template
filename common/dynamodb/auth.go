package dynamodb

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Auth struct {
	PK string `json:"PK"` // uid#
	SK string `json:"SK"` // auth#

	Token         interface{} `json:"token,omitempty"` // 인즈 토큰
	SessionKey    string      `json:"session_key"`     // 세션키
	SessionKeyTTL int64       `json:"session_key_ttl"` // 세션키 TTL
	Attr          interface{} `json:"attr,omitempty"`  // 유저 속성

	Grade string `json:"grade"` // 유저 등급	("" 일반 유저, "ADMIN", "SUPER_USER")
	Pw    string `json:"pw"`    // 유저 암호

	UpdateAt int64 `json:"update_at"`
	CreateAt int64 `json:"create_at"`
	TTL      int64 `json:"ttl"`
}

type AdminAuth struct {
	Id         string `json:"id"`
	SessionKey string `json:"session_key"`
	Grade      string `json:"grade"`
}

var AuthTableTTL = int64(3600 * 24 * 365) // 1 year
var SessionKeyTTL = int64(3600 * 8)       // 8 hour

// GetItem ...
func GetAuth(ctx context.Context, pk string) (*Auth, error) {
	query := `SELECT * FROM "game-oz" WHERE "PK" = ? AND "SK" = 'AUTH#'`
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

	auth := &Auth{}
	if len(result.Items) > 0 {
		UnmarshalMap(result.Items[0], auth)
		return auth, nil
	}
	return nil, errors.New("no account data")
}

func GetAuthBySessionKey(ctx context.Context, ssk string) (*Auth, error) {
	query := `SELECT * FROM "game-oz" WHERE "session_key" = ? AND "SK" = 'AUTH#'`
	param := dynamodb.ExecuteStatementInput{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: ssk},
		},
	}
	result, err := dynamodbconn.ExecuteStatement(ctx, &param)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	auth := &Auth{}
	if len(result.Items) > 0 {
		UnmarshalMap(result.Items[0], auth)
		return auth, nil
	}
	return nil, errors.New("no account data")
}

func GetAuthList(ctx context.Context) ([]*Auth, error) {
	query := `
		SELECT PK, session_key, grade 
		FROM "game-oz"."session_key-SK-index"
		WHERE "SK" = 'AUTH#'
		AND "grade" = 'ADMIN'
	`

	param := dynamodb.ExecuteStatementInput{
		Statement: aws.String(query),
	}

	result, err := dynamodbconn.ExecuteStatement(ctx, &param)
	if err != nil {
		log.Println(err)
	}
	ret := make([]*Auth, 0)
	for _, v := range result.Items {
		item := Auth{}
		err = UnmarshalMap(v, &item)
		if err != nil {
			log.Println(err)
			continue
		}
		ret = append(ret, &item)
	}
	return ret, nil
}

// PutItem ...
func PutAuth(ctx context.Context, auth *Auth) error {
	query := `INSERT INTO "game-oz"
	VALUE {
		'PK':?, 
		'SK':'AUTH#',
		'token':?,
		'session_key':?,
		'session_key_ttl':?,
		'attr':?,
		'grade':?,
		'pw':?,
		'update_at':?,
		'create_at':?,
		'ttl':?
	}`
	at := time.Now().Unix()
	token, _ := MarshalMap(auth.Token)
	attr, _ := MarshalMap(auth.Attr)

	param := dynamodb.ExecuteStatementInput{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: auth.PK},
			&types.AttributeValueMemberM{Value: token},
			&types.AttributeValueMemberS{Value: auth.SessionKey},
			&types.AttributeValueMemberN{Value: IntToStr(at + SessionKeyTTL)},
			&types.AttributeValueMemberM{Value: attr},
			&types.AttributeValueMemberS{Value: auth.Grade},
			&types.AttributeValueMemberS{Value: auth.Pw},
			&types.AttributeValueMemberN{Value: IntToStr(at)},
			&types.AttributeValueMemberN{Value: IntToStr(at)},
			&types.AttributeValueMemberN{Value: IntToStr(at + AuthTableTTL)},
		},
	}
	_, err := dynamodbconn.ExecuteStatement(ctx, &param)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// UpdateItem ...
func UpdateAuth(ctx context.Context, auth *Auth) error {
	query := `UPDATE "game-oz"
			SET "token" = ?,
				"attr" = ?,
				"session_key" = ?,
				"grade" = ?,
				"session_key_ttl" = ?,
				"update_at" = ?,
				"ttl" = ?
			WHERE "PK" = ? AND "SK" = 'AUTH#'`

	at := time.Now().Unix()
	token, _ := MarshalMap(auth.Token)
	attr, _ := MarshalMap(auth.Attr)

	param := dynamodb.ExecuteStatementInput{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberM{Value: token},
			&types.AttributeValueMemberM{Value: attr},
			&types.AttributeValueMemberS{Value: auth.SessionKey},
			&types.AttributeValueMemberS{Value: auth.Grade},
			&types.AttributeValueMemberN{Value: IntToStr(at + SessionKeyTTL)},
			&types.AttributeValueMemberN{Value: IntToStr(at)},
			&types.AttributeValueMemberN{Value: IntToStr(at + AuthTableTTL)},
			&types.AttributeValueMemberS{Value: auth.PK},
		},
	}
	_, err := dynamodbconn.ExecuteStatement(context.Background(), &param)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteItem ...
func DeleteAuth(ctx context.Context, pk string) error {
	query := `DELETE FROM "game-oz"			
			WHERE "PK" = ? AND "SK" = 'AUTH#'
			AND "grade" = 'ADMIN'`

	param := dynamodb.ExecuteStatementInput{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: pk},
		},
	}
	_, err := dynamodbconn.ExecuteStatement(context.Background(), &param)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
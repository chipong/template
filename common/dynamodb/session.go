package dynamodb

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

func SessionKeyToUidHeader(c *gin.Context) {
	if strings.Contains(c.Request.RequestURI, "simulation") {
		return
	}
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization",
		})
		c.Abort()
		return
	}

	prefix := "Bearer "
	token := authorization

	if strings.HasPrefix(authorization, prefix) {
		token = authorization[len(prefix):]
	}

	query := `SELECT "PK" FROM "game-oz"."session_key-index" WHERE "session_key" = ?`
	param := dynamodb.ExecuteStatementInput{
		Statement: aws.String(query),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: token},
		},
	}
	result, err := dynamodbconn.ExecuteStatement(c.Request.Context(), &param)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": err,
		})
		c.Abort()
		return
	}

	if len(result.Items) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization key",
		})
		c.Abort()
		return
	}

	uid := result.Items[0]["PK"].(*types.AttributeValueMemberS).Value
	c.Request.Header.Set("UID", uid)
	//c.Next()
}

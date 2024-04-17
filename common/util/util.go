package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/aws/aws-lambda-go/events"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/chipong/template/common/dynamodb"
	"github.com/chipong/template/common/redisCache"
)

func GetHeaderUidShard(c *gin.Context) (string, int, error) {
	uid := c.Request.Header.Get("UID")
	ShardIndex := c.Request.Header.Get("SHARD_INDEX")
	shardIndex, err := strconv.Atoi(ShardIndex)
	if err != nil {
		return uid, 0, err
	}
	return uid, shardIndex, nil
}

func GetHeaderUid(c *gin.Context) (string, error) {
	uid := c.Request.Header.Get("UID")
	if uid == "" {
		return uid, errors.New("no have uid")
	}
	return uid, nil
}

// Bearer를 위해 수정
func GetHeaderUidSessionKey(c *gin.Context) (string, string, error) {
	uid := c.Request.Header.Get("UID")
	authorization := c.Request.Header.Get("Authorization")
	var sessionKey string

	if uid == "" {
		uid = "not exist"
	}

	if authorization == "" {
		sessionKey = "not exist"
	} else {
		prefix := "Bearer "
		sessionKey = authorization

		if strings.HasPrefix(authorization, prefix) {
			sessionKey = authorization[len(prefix):]
		}
	}

	return uid, sessionKey, nil
}

func GetHeaderSessionKey(c *gin.Context) string {
	temp := c.Request.Header.Get("Authorization")
	splitStr := strings.Split(temp, " ")
	if len(splitStr) < 2 {
		return "not exist"
	}

	return splitStr[1]
}

func GetAuthorizedUIDSSK(c *gin.Context) (string, string, error) {
	ssKey := c.Request.Header.Get("Authorization")
	if ssKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization",
		})
		c.Abort()
		return "", "", errors.New("No authorization")
	}

	prefix := "Bearer "
	token := ssKey

	log.Println("Authorization: ", ssKey)

	if strings.HasPrefix(ssKey, prefix) {
		token = ssKey[len(prefix):]
	}

	auth, err := dynamodb.GetAuthBySessionKey(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization",
		})
		c.Abort()
		return "", "", err
	}

	pkSplit := strings.Split(auth.PK, "#")
	log.Println(pkSplit)

	var uid string
	if len(pkSplit) > 1 {
		uid = pkSplit[1]
	} else {
		uid = pkSplit[0]
	}

	log.Println("uid: ", uid, "sessionKey: ", token)
	c.Request.Header.Set("UID", uid)
	return uid, token, nil
}

func ShardIndex(c *gin.Context) (int, error) {
	shardIndex := c.Request.Header.Get("SHARD_INDEX")
	//log.Println(c.Request.Header)
	return strconv.Atoi(shardIndex)
}

func RequestHttp(url, method string, body io.Reader) ([]byte, error) {
	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		ans := &struct {
			ErrCode int    `json:"err_code"`
			ErrMsg  string `json:"err_msg"`
		}{}
		err = json.Unmarshal(respBody, ans)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ans.ErrMsg)
	}
	return respBody, nil
}

func RequestHttpWithContext(c *gin.Context, uid string, url, method string, body io.Reader) ([]byte, error) {
	defer func() {
		if cover := recover(); cover != nil {
			log.Println(c, "RequestHttpWithContexterr: ", cover)
			return
		}
	}()
	
	log.Println(url)
	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	httpReq.Header.Add("Content-Type", "application/json")
	if uid != "" {
		httpReq.Header.Set("UID", uid)
	}

	httpReq.Header.Add("Authorization", c.Request.Header.Get("Authorization"))
	httpReq.Header.Add("User-Agent", c.Request.Header.Get("User-Agent"))
	httpReq.Header.Add("X-Forwarded-For", c.ClientIP())

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		ans := &struct {
			ErrCode int    `json:"err_code"`
			ErrMsg  string `json:"err_msg"`
		}{}
		err = json.Unmarshal(respBody, ans)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ans.ErrMsg)
	}

	return respBody, nil
}

func Hash(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	md := h.Sum(nil)
	dest := hex.EncodeToString(md)
	return dest
}

func Unmarshal(r io.Reader, m proto.Message) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err := protojson.Unmarshal(body, m); err != nil {
		return err
	}
	return nil
}

func GetHeaderUidWebsocket(event events.APIGatewayWebsocketProxyRequest) (string, error) {
	authorization := event.Headers["Authorization"]
	if authorization == "" {
		return "", errors.New("no authorization")
	}

	return redisCache.GetSessionKey(authorization)
}

func GoroutineID() string {
	gr := bytes.Fields(debug.Stack())[1]
	return string(gr)
}

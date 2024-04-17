package redisCache

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

var sessionTTL time.Duration = time.Hour * 8

func SetSessionTTL(ttl time.Duration) {
	sessionTTL = ttl
}

func IsConnected() error {
	if read == nil || write == nil {
		return errors.New("redis client no connection")
	}
	return nil
}

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

	key := fmt.Sprintf("%s:%s:session", appName, token)

	// redis key check
	uid := ""
	var err error
	if cacheCfg.Ver != "" {
		uid, err = read.GetEx(context.Background(), key, sessionTTL).Result()
	} else {
		uid, err = read.Get(context.Background(), key).Result()
	}
	if err != nil {
		log.Println(err, " ", token)
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis no authorization key: " + token,
		})
		c.Abort()
		return
	}

	if uid == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization key",
		})
		c.Abort()
		return
	}

	if cacheCfg.Ver == "" {
		Touch("session", uid)
	}

	c.Request.Header.Set("UID", uid)
	//c.Next()
}

func DelSessionKey(sessionKey string) error {
	if write == nil {
		return errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:%s:session", appName, sessionKey)
	_, err := write.Del(context.Background(), key).Result()
	if err != nil {

		return err
	}
	return nil
}

func SetSessionKey(uid, sessionKey string) error {
	if write == nil {
		return errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:%s:session", appName, sessionKey)
	_, err := write.Set(context.Background(), key, uid, sessionTTL).Result()
	if err != nil {

		return err
	}
	return nil
}

func GetSessionKey(authorization string) (string, error) {
	prefix := "Bearer "
	token := authorization

	if strings.HasPrefix(authorization, prefix) {
		token = authorization[len(prefix):]
	}

	key := fmt.Sprintf("%s:%s:session", appName, token)

	// redis key check
	uid, err := read.Get(context.Background(), key).Result()
	if err != nil {
		log.Println(err, " ", token)
		return "", errors.New("no authorization key")
	}

	if uid == "" {
		return "", errors.New("no authorization key")
	}

	return uid, nil
}

package core

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/redis.v4"
)

// WriteRedisClient ...
var WriteRedisClient *redis.Client

// ReadRedisClient ...
var ReadRedisClient *redis.Client

// SessionCheckConfig ...
type SessionCheckConfig struct {
	PAddr    string `yaml:"primary_addr"`
	RAddr    string `yaml:"reader_addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Expire   int    `yaml:"expire"`
	AppName  string `yaml:"app_name"`
}

var sessionCfg = &SessionCheckConfig{}

// InitSessionCheck ...
func InitSessionCheckEx(cfg SessionCheckConfig) {
	// redis database 1 connecttion
	WriteRedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.PAddr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	if WriteRedisClient == nil {
		log.Println("redis connect fail")
		return
	}

	// redis database 1 connecttion
	ReadRedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RAddr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	if ReadRedisClient == nil {
		log.Println("redis connect fail")
		return
	}

	sessionCfg = &cfg

	pingRedis()
}

func InitSessionCheck(cfg SessionCheckConfig) {
	// redis database 1 connecttion
	WriteRedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.PAddr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	if WriteRedisClient == nil {
		log.Println("redis connect fail")
		return
	}

	// redis database 1 connecttion
	ReadRedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.RAddr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	if ReadRedisClient == nil {
		log.Println("redis connect fail")
		return
	}

	sessionCfg = &cfg

	pingRedis()

	// config reload
	go func() {
		c := time.Tick(time.Minute)
		for range c {
			//log.Println("redis ping check: ", now.UTC())
			_, err := WriteRedisClient.Ping().Result()
			if err != nil {
				log.Println("redis master connection fail")
				WriteRedisClient = redis.NewClient(&redis.Options{
					Addr:     sessionCfg.PAddr,
					Password: sessionCfg.Password, // no password set
					DB:       sessionCfg.DB,       // use default DB
				})
			}

			_, err = ReadRedisClient.Ping().Result()
			if err != nil {
				log.Println("redis slave connection fail")
				ReadRedisClient = redis.NewClient(&redis.Options{
					Addr:     sessionCfg.RAddr,
					Password: sessionCfg.Password, // no password set
					DB:       sessionCfg.DB,       // use default DB
				})
			}
		}
	}()
}

func pingRedis() {
	_, err := WriteRedisClient.Ping().Result()
	if err != nil {
		log.Println("redis master connection fail")
		return
	}

	_, err = ReadRedisClient.Ping().Result()
	if err != nil {
		log.Println("redis slave connection fail")
		return
	}

	log.Println("redis connection success")
}

func Ping() {
	_, err := WriteRedisClient.Ping().Result()
	if err != nil {
		log.Println("redis master connection fail")
		WriteRedisClient = redis.NewClient(&redis.Options{
			Addr:     sessionCfg.PAddr,
			Password: sessionCfg.Password, // no password set
			DB:       sessionCfg.DB,       // use default DB
		})
	}

	_, err = ReadRedisClient.Ping().Result()
	if err != nil {
		log.Println("redis slave connection fail")
		ReadRedisClient = redis.NewClient(&redis.Options{
			Addr:     sessionCfg.RAddr,
			Password: sessionCfg.Password, // no password set
			DB:       sessionCfg.DB,       // use default DB
		})
	}
}

// SessionCheck ..
func SessionCheck(c *gin.Context) {
	ssKey := c.Request.Header.Get("Authorization")
	if ssKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization",
		})
		c.Abort()
		return
	}

	prefix := "Bearer "
	token := ssKey

	if strings.HasPrefix(ssKey, prefix) {
		token = ssKey[len(prefix):]
	}

	if ReadRedisClient == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis client no connection",
		})
		c.Abort()
		return
	}

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, token)

	// redis key check
	uid, err := ReadRedisClient.Get(key).Result()
	if err != nil {
		log.Println(err, " ", token)
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis no authorization key: " + token,
		})
		c.Abort()
		return
	}

	if uid == "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "authorization key is admin",
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

	SessionKeyTTL(token)

	c.Request.Header.Set("UID", uid)

	c.Next()
}

// SessionCheck ..
func SessionCheckAndShardIndex(c *gin.Context) {
	ssKey := c.Request.Header.Get("Authorization")
	if ssKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization",
		})
		c.Abort()
		return
	}

	prefix := "Bearer "
	token := ssKey

	if strings.HasPrefix(ssKey, prefix) {
		token = ssKey[len(prefix):]
	}

	if ReadRedisClient == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis client no connection",
		})
		c.Abort()
		return
	}

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, token)

	// redis key check
	uid, err := ReadRedisClient.Get(key).Result()
	if err != nil {
		log.Println(err, " ", token)
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis no authorization key: " + token,
		})
		c.Abort()
		return
	}

	if uid == "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "authorization key is admin",
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

	// redis key check
	shardKey := fmt.Sprintf("%s:shard:%s", sessionCfg.AppName, uid)
	shardIndex, err := ReadRedisClient.Get(shardKey).Result()
	if err != nil {
		log.Println(err, " ", uid)
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis no shard key: " + uid,
		})
		c.Abort()
		return
	}

	SessionKeyTTL(token)

	c.Request.Header.Set("UID", uid)
	c.Request.Header.Set("SHARD_INDEX", shardIndex)

	c.Next()
}

func SessionCheckAndShardIndexProcCheck(c *gin.Context) {
	ssKey := c.Request.Header.Get("Authorization")
	if ssKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization",
		})
		c.Abort()
		return
	}

	prefix := "Bearer "
	token := ssKey

	if strings.HasPrefix(ssKey, prefix) {
		token = ssKey[len(prefix):]
	}

	if ReadRedisClient == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis client no connection",
		})
		c.Abort()
		return
	}

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, token)

	// redis key check
	uid, err := ReadRedisClient.Get(key).Result()
	if err != nil {
		log.Println(err, " ", token)
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis no authorization key: " + token,
		})
		c.Abort()
		return
	}

	if uid == "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "authorization key is admin",
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

	// redis key check
	shardKey := fmt.Sprintf("%s:shard:%s", sessionCfg.AppName, uid)
	shardIndex, err := ReadRedisClient.Get(shardKey).Result()
	if err != nil {
		log.Println(err, " ", uid)
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis no shard key: " + uid,
		})
		c.Abort()
		return
	}

	SessionKeyTTL(token)

	c.Request.Header.Set("UID", uid)
	c.Request.Header.Set("SHARD_INDEX", shardIndex)

	c.Next()
}

// SessionCheckForAdmin ...
func SessionCheckForAdmin(c *gin.Context) {
	ssKey := c.Request.Header.Get("Authorization")
	if ssKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization",
		})
		c.Abort()
		return
	}

	prefix := "Bearer "
	token := ssKey

	if strings.HasPrefix(ssKey, prefix) {
		token = ssKey[len(prefix):]
	}

	if ReadRedisClient == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis client no connection",
		})
		c.Abort()
		return
	}

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, token)
	// redis key check
	uid, err := ReadRedisClient.Get(key).Result()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "redis no authorization key: " + token,
		})
		c.Abort()
		return
	}

	if uid != "admin" {

		c.JSON(http.StatusUnauthorized, gin.H{
			"errMsg": "no authorization key",
		})
		c.Abort()
		return
	}

	c.Request.Header.Set("UID", uid)

	c.Next()
}

// GenSessionKey ...
func GenSessionKey(id string, expire time.Duration) (string, error) {
	if WriteRedisClient == nil {
		return "", errors.New("redis client no connection")
	}

	sessionKey := GenUUID()

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, sessionKey)
	_, err := WriteRedisClient.Set(key, id, expire).Result()
	if err != nil {

		return "", err
	}
	return sessionKey, nil
}

// SetSessionKey ...
func SetSessionKey(id string) (string, error) {
	if WriteRedisClient == nil {
		return "", errors.New("redis client no connection")
	}

	sessionKey := GenUUID()
	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, sessionKey)
	_, err := WriteRedisClient.Set(key, id, time.Second*3600*time.Duration(sessionCfg.Expire)).Result()
	if err != nil {

		return "", err
	}
	return sessionKey, nil
}

// SetSessionKey ...
func SetSessionKeyNoKeyGen(id, sskey string) error {
	if WriteRedisClient == nil {
		return errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, sskey)
	_, err := WriteRedisClient.Set(key, id, time.Second*3600*time.Duration(sessionCfg.Expire)).Result()
	if err != nil {

		return err
	}
	return nil
}

// DelSessionKey ...
func DelSessionKey(sessionKey string) error {
	if WriteRedisClient == nil {
		return errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, sessionKey)
	_, err := WriteRedisClient.Del(key).Result()
	if err != nil {

		return err
	}
	return nil
}

// SessionKeyCount ...
func SessionKeyCount() (int, error) {
	if ReadRedisClient == nil {
		return 0, errors.New("redis client no connection")
	}
	keys, _, err := ReadRedisClient.Scan(0, sessionCfg.AppName+":session:*", 10000).Result()
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

// SessionKeyTTL ..
func SessionKeyTTL(token string) error {
	if WriteRedisClient == nil {
		return errors.New("redis client no connection")
	}

	if sessionCfg.Expire == 0 {
		return nil
	}

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, token)
	_, err := WriteRedisClient.Expire(key, time.Second*3600*time.Duration(sessionCfg.Expire)).Result()
	if err != nil {

		return err
	}
	return nil
}

func GetSessionKeyAndShardIndex(authorization string) (string, int, error) {
	prefix := "Bearer "
	token := authorization

	if strings.HasPrefix(authorization, prefix) {
		token = authorization[len(prefix):]
	}

	if ReadRedisClient == nil {
		return "", 0, errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:session:%s", sessionCfg.AppName, token)

	// redis key check
	uid, err := ReadRedisClient.Get(key).Result()
	if err != nil {
		log.Println(err, " ", token)
		return "", 0, errors.New("no authorization key")
	}

	if uid == "" {
		return "", 0, errors.New("no authorization key")
	}

	// redis key check
	shardKey := fmt.Sprintf("%s:shard:%s", sessionCfg.AppName, uid)
	shardIndex, err := ReadRedisClient.Get(shardKey).Result()
	if err != nil {
		log.Println(err, " ", uid)
		return "", 0, errors.New("no shard key")
	}

	shard, _ := ParseInt(shardIndex)
	return uid, shard, nil
}

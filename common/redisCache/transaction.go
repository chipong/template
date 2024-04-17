package redisCache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	// TransactionTimeOut ...
	TransactionTimeOut = time.Second * 300
)

func hash(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	md := h.Sum(nil)
	dest := hex.EncodeToString(md)
	return dest
}

// func hash2(src string) string {
// 	dest := md5.Sum([]byte(src))
// 	return string(dest[:])
// }

func BeginTran(uid string, cn, query string) (string, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		// 50ms
		if elapsed > time.Millisecond*50 {
			log.Println("BeginTran took: ", elapsed)
		}
	}()

	key := fmt.Sprintf("%s:transaction:%s:%s", appName, hash(cn), uid)
	transKey := hash(key)
	_, err := write.Set(context.Background(), key, query, TransactionTimeOut).Result()
	if err != nil {
		return "", err
	}

	_, err = write.Set(context.Background(), transKey, cn, TransactionTimeOut).Result()
	if err != nil {
		return "", err
	}
	return transKey, nil
}

func GetCn(transkey string) (string, error) {
	return read.Get(context.Background(), transkey).Result()
}

func CheckTranQuery(uid string, cn string) error {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		// 50ms
		if elapsed > time.Millisecond*50 {
			log.Println("CheckTranQuery took: ", elapsed)
		}
	}()

	key := fmt.Sprintf("%s:transaction:%s:%s", appName, hash(cn), uid)
	data, err := read.Get(context.Background(), key).Result()
	if err == nil && data != "" {
		log.Println("already key: ", key)
		log.Println("already tx: ", data)
		return errors.New("already transaction beginning")
	}

	return nil
}

func GetTranQuery(uid string, cn, transKey string) (string, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		// 50ms
		if elapsed > time.Millisecond*50 {
			log.Println("GetTranQuery took: ", elapsed)
		}
	}()
	key := fmt.Sprintf("%s:transaction:%s:%s", appName, hash(cn), uid)

	if hash(key) != transKey {
		return "", errors.New("invalid transaction key")
	}

	data, err := read.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	return data, nil
}

func CompleteTran(uid string, cn, transKey string) error {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		// 50ms
		if elapsed > time.Millisecond*50 {
			log.Println("CompleteTran took: ", elapsed)
		}
	}()

	key := fmt.Sprintf("%s:transaction:%s:%s", appName, hash(cn), uid)

	if hash(key) != transKey {
		log.Println(key)
		log.Println(hash(key))
		log.Println(transKey)
		return errors.New("invalid transaction key")
	}

	write.Del(context.Background(), transKey).Result()
	_, err := write.Del(context.Background(), key).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

const (
	// ProcTimeOut ...
	ProcTimeOut = time.Second * 5
)

func BeginProc(uid string, cn string) error {
	if err := IsConnected(); err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s:proc", appName, uid)
	data, err := read.Get(context.Background(), key).Result()
	if data != "" {
		return errors.New("already processing")
	}

	_, err = write.Set(context.Background(), key, cn, ProcTimeOut).Result()
	if err != nil {
		return err
	}
	return nil
}

func EndProc(uid string) error {
	if err := IsConnected(); err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s:proc", appName, uid)

	_, err := write.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}

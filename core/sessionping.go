package core

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	// PingTTL ...
	PingTTL = time.Minute * 10
)

// Caching ...
func Caching(appName, sskey, attr string) error {
	if WriteRedisClient == nil {
		return errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:auth:%s", appName, sskey)
	data, err := WriteRedisClient.Get(key).Result()
	if data != "" {
		return Touch(appName, sskey)
	}

	_, err = WriteRedisClient.Set(key, attr, PingTTL).Result()
	{
		ccukey := fmt.Sprintf("%s:monitor:ccu-count", appName)
		WriteRedisClient.Incr(ccukey).Result()
	}
	return err
}

func Uncaching(appName, sskey string) error {
	if WriteRedisClient == nil {
		return errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:auth:%s", appName, sskey)
	_, err := WriteRedisClient.Del(key).Result()
	return err
}

// CCUFullScan ...
func CCUFullScan(appName string) (int, error) {
	if ReadRedisClient == nil {
		return 0, errors.New("redis client no connection")
	}

	keys, _, err := ReadRedisClient.Scan(0, appName+":auth:*", 100000).Result()
	if err != nil {
		return 0, err
	}
	ccu := len(keys)

	// update ccu count
	if WriteRedisClient != nil {
		key := fmt.Sprintf("%s:monitor:ccu-count", appName)
		WriteRedisClient.Set(key, strconv.Itoa(ccu), PingTTL).Result()
	}
	return ccu, nil
}

// CCU ...
func CCU(appName string) (int, error) {
	if ReadRedisClient == nil {
		return 0, errors.New("redis client no connection")
	}

	r, err := ReadRedisClient.Get(appName + ":monitor:ccu-count:*").Result()
	if err != nil {
		return 0, err
	}
	ccu, _ := strconv.Atoi(r)
	return ccu, nil
}

// Touch ...
func Touch(appName, sskey string) error {
	if WriteRedisClient == nil {
		return errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:auth:%s", appName, sskey)
	_, err := WriteRedisClient.Expire(key, PingTTL).Result()
	return err
}

// GetCache ...
func GetCache(appName, sskey string) (string, error) {
	if ReadRedisClient == nil {
		return "", errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:auth:%s", appName, sskey)
	r, err := ReadRedisClient.Get(key).Result()
	if err != nil {
		return "", err
	}

	return r, nil
}

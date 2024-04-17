package redisCache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"strconv"

	redis "github.com/go-redis/redis/v8"
	"github.com/chipong/template/common/proto"
)

var (
	write    *redis.Client
	read     *redis.Client
	cacheCfg RedisCacheConfig
)

type RedisCacheConfig struct {
	PAddr    string `yaml:"primary_addr"`
	RAddr    string `yaml:"reader_addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	Expire   int    `yaml:"expire"`
	Ver      string `yaml:"ver"`
}

const (
	// PingTTL ...
	PingTTL       = time.Minute * 10
	DataTTL       = time.Hour * 3
	appName       = "macovill.oz"
	pageSize      = 10
	pageListCount = 10
)

func Initialize(cfg RedisCacheConfig) {
	cacheCfg = cfg

	write = redis.NewClient(&redis.Options{
		Addr:     cfg.PAddr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	read = redis.NewClient(&redis.Options{
		Addr:     cfg.RAddr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	go pingRedis(write, read)

	// config reload
	go func() {
		c := time.Tick(time.Minute)
		for range c {

			_, err := write.Ping(context.Background()).Result()
			if err != nil {
				write = redis.NewClient(&redis.Options{
					Addr:     cacheCfg.PAddr,
					Password: cacheCfg.Password, // no password set
					DB:       cacheCfg.DB,       // use default DB
				})

				if write != nil {
					log.Println("redis master reconnection")
				}
			}

			_, err = read.Ping(context.Background()).Result()
			if err != nil {
				log.Println("redis slave connection fail")
				read = redis.NewClient(&redis.Options{
					Addr:     cacheCfg.RAddr,
					Password: cacheCfg.Password, // no password set
					DB:       cacheCfg.DB,       // use default DB
				})

				if read != nil {
					log.Println("redis master reconnection")
				}
			}
		}
	}()
}

func InitializeEx(cfg RedisCacheConfig) {
	cacheCfg = cfg
	write = redis.NewClient(&redis.Options{
		Addr:     cfg.PAddr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	read = redis.NewClient(&redis.Options{
		Addr:     cfg.RAddr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	SetSessionTTL(time.Hour * time.Duration(cacheCfg.Expire))

	pingRedis(write, read)
	log.Println("cache connection - ", cfg.RAddr)
}

func pingRedis(w, r *redis.Client) {
	_, err := w.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis master connection fail")
		return
	}

	_, err = r.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis slave connection fail")
		return
	}
}

func Ping() {
	_, err := write.Ping(context.Background()).Result()
	if err != nil {
		write = redis.NewClient(&redis.Options{
			Addr:     cacheCfg.PAddr,
			Password: cacheCfg.Password, // no password set
			DB:       cacheCfg.DB,       // use default DB
		})
	}

	_, err = read.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis slave connection fail")
		read = redis.NewClient(&redis.Options{
			Addr:     cacheCfg.RAddr,
			Password: cacheCfg.Password, // no password set
			DB:       cacheCfg.DB,       // use default DB
		})
	}
}

// Touch ...
func Touch(contents, uid string) error {
	if write == nil {
		return errors.New("redis client no connection")
	}

	key := fmt.Sprintf("%s:%s:%s", appName, uid, contents)
	_, err := write.Expire(context.Background(), key, PingTTL).Result()
	return err
}

func Set(key, value string, ttl time.Duration) error {
	if write == nil {
		return errors.New("redis client no connection")
	}

	_, err := write.Set(context.Background(), key, value, ttl).Result()
	if err != nil {
		return err
	}
	return nil
}

func Del(key string) error {
	if write == nil {
		return errors.New("redis client no connection")
	}

	_, err := write.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}

func Scan(cursor uint64, key string, count int64) ([]string, error) {
	if read == nil {
		return nil, errors.New("redis client no connection")
	}

	keys, _, err := read.Scan(context.Background(), cursor, key, count).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func Eval(script string, keys, args []string) error {
	if write == nil {
		return errors.New("redis client no connection")
	}
	out := write.Eval(context.Background(), script, keys, args)
	if out.Err() != nil {
		return out.Err()
	}
	//log.Println(out.String())
	return nil
}

func MGet[T oz.OZTemplate ](
	key string, ttl time.Duration, dummy *T) ([]*T, error) {

	ttl = ttl / time.Second
	script := `
	local out = {};
	local v = redis.call('KEYS', KEYS[1]);
	for index, key in ipairs(v) do	
		out[index] = redis.call('GET', key);
		redis.call('EXPIRE', key, ARGV[1]);
	end 
	return out;`

	results, err := write.Eval(
		context.Background(),
		script,
		[]string{key},
		[]string{strconv.FormatInt(int64(ttl), 10)}).Result()
	if err != nil {
		return nil, err
	}

	outs := make([]*T, 0)
	for _, v := range results.([]interface{}) {
		temp := []*T{}
		err = json.Unmarshal([]byte(v.(string)), &temp)
		if err != nil {
			break
		}
		outs = append(outs, temp...)
	}
	return outs, nil
}

func MExpire(key string, ttl time.Duration) error {

	ttl = ttl / time.Second
	script := `local v = redis.call('KEYS', KEYS[1]);for index, key in ipairs(v) do	redis.call('EXPIRE', key, ARGV[1]); end return 'OK';`
	err := Eval(script, []string{key}, []string{strconv.FormatInt(int64(ttl), 10)})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func DelKeys(keys []string) error {
	script := `
		local i = 1
		for k, v in ipairs(KEYS) do
			redis.call('DEL', v);
			i = i + 1;
		end 
		return 'OK';`
	out := Eval(script, keys, []string{})
	log.Println(out)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	return nil
}

func DelPattern(uid, key string) error {
	keys := fmt.Sprintf("%s:%s:%s*", appName, uid, key)
	script := `
		local v = redis.call('KEYS', KEYS[1]);
		for index, key in ipairs(v) do	
			redis.call('DEL', key);
		end 
		return 'OK';`
	err := Eval(script, []string{keys}, []string{})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func LDelAndPush[T oz.OZTemplate](
	key string, ttl time.Duration, datas []*T) error {
	values := []string{}
	for _, v := range datas {
		jsonData, err := json.Marshal(v)
		if err != nil {
			return err
		}
		values = append(values, string(jsonData))
	}

	ttl = ttl / time.Second
	script := fmt.Sprintf(`
		redis.call('DEL', KEYS[1]);
		for k, v in ipairs(ARGV) do
			redis.call('LPUSH', KEYS[1], v);
		end
		redis.call('EXPIRE', KEYS[1], %d );
		return 'OK';`, ttl)
	err := Eval(script, []string{key}, values)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func LRange[T oz.OZTemplate ](
	key string, ttl time.Duration, datas *T) ([]*T, error) {

	ttl = ttl / time.Second
	script := fmt.Sprintf(`
		local out = redis.call('LRANGE', KEYS[1], 0, -1);
		redis.call('EXPIRE', KEYS[1], %d );
		return out;`, ttl)
	result, err := write.Eval(context.Background(), script, []string{key}, []string{}).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	outs := make([]*T, 0)
	for _, v := range result.([]interface{}) {
		temp := new(T)
		err := json.Unmarshal([]byte(v.(string)), temp)
		if err != nil {
			break
		}
		outs = append(outs, temp)
	}
	return outs, nil
}

// score 갱신 후 현재 rank return
func ZAdd(key string, ttl time.Duration, uid string, score int64) (int64, error) {
	ttl = ttl / time.Second
	script := fmt.Sprintf(`
		local out = redis.call('ZADD', KEYS[1], ARGV[1], ARGV[2]);
		--redis.call('EXPIRE', KEYS[1], %d);
		return out;
	`, ttl)
	_, err := write.Eval(context.Background(), script, []string{key}, 
	[]string{strconv.FormatInt(int64(score), 10), uid}).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	script = fmt.Sprintln(`
		local out = redis.call('ZRANK', KEYS[1], ARGV[1]);
		return out;`)
	rank, err := write.Eval(context.Background(), script, []string{key}, []string{uid}).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	
	return (rank.(int64) + 1), nil
}

/*
// zincrby 해당 값 만큼 기존 값에 더한다.
func ZincrBy(key string, ttl time.Duration, uid string, inc_score int64, at int32) (int64, error) {
	ttl = ttl / time.Second
	script := fmt.Sprintf(`
		local out = redis.call('ZINCRBY', KEYS[1], ARGV[1], ARGV[2]);
		redis.call('EXPIRE', KEYS[1], %d);
		return out;
	`, ttl)
	_, err := write.Eval(context.Background(), script, []string{key}, 
	[]string{strconv.FormatInt(int64(inc_score * -1), 10), uid}).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	script = fmt.Sprintln(`
		local out = redis.call('ZRANK', KEYS[1], ARGV[1]);
		return out;`)
	rank, err := write.Eval(context.Background(), script, []string{key}, []string{uid}).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	
	return (rank.(int64) + 1), nil
}
*/

func ZRem(key string, uid string) (int64, error) {
	script := fmt.Sprintln(`
		local out = redis.call('ZREM', KEYS[1], ARGV[1]);
		return out;
	`)
	result, err := write.Eval(context.Background(), script, []string{key}, []string{uid}).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	
	return result.(int64), nil
}

func ZRank(key string, ttl time.Duration, uid string) (int64, error) {
	ttl = ttl / time.Second
	script := fmt.Sprintf(`
		local out = redis.call('ZRANK', KEYS[1], ARGV[1]);
		--redis.call('EXPIRE', KEYS[1], %d);
		return out;`, ttl)
	rank, err := write.Eval(context.Background(), script, []string{key}, []string{uid}).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	
	return (rank.(int64) + 1), nil
}

func ZRange(key string, ttl time.Duration, start, end int64) (interface{}, error) {
	ttl = ttl / time.Second
	script := fmt.Sprintf(`
		local out = redis.call('ZRANGE', KEYS[1], ARGV[1], ARGV[2], 'withscores');
		--redis.call('EXPIRE', KEYS[1], %d);
		return out;`, ttl)
		
	results, err := write.Eval(context.Background(), script, []string{key}, 
		[]string{strconv.FormatInt(int64(start), 10), strconv.FormatInt(int64(end), 10)}).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return results, nil
}
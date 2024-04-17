package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/signal"
	"syscall"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/gzip"

	"github.com/chipong/template/template-server/api/template"
	"github.com/chipong/template/template-server/router"

	"github.com/chipong/template/core"
	awss3 "github.com/chipong/template/common/awsS3"
	awsssm "github.com/chipong/template/common/awsSSM"
	"github.com/chipong/template/common/dynamodb"
	"github.com/chipong/template/common/mvlog"
	"github.com/chipong/template/common/redisCache"
	"github.com/chipong/template/common/slack"
)

const (
	appName = "template"
	serverName = "template-server"
	configName = "template-server-config"
)

// Config ...
type Config struct {
	// kinesis
	Kinesis struct {
		Region     string `yaml:"region"`
		AccessKey  string `yaml:"access_key"`
		SecretKey  string `yaml:"secret_access_key"`
		StreamName string `yaml:"stream_name"`
	} `yaml:"kinesis,omitempty"`

	S3 struct {
		Region    string `yaml:"region"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_access_key"`
		Bucket    string `yaml:"bucket"`
	} `yaml:"s3,omitempty"`

	DynamoDB struct {
		Region    string `yaml:"region"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_access_key"`
		EndPoint  string `yaml:"end_point"`
	} `yaml:"dynamodb,omitempty"`

	Cache struct {
		PAddr    string `yaml:"primary_addr"`
		RAddr    string `yaml:"reader_addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		Expire   int    `yaml:"expire"`
		Ver      string `yaml:"ver"`
	} `yaml:"cache,omitempty"`

	Slack struct {
		Token 		string `yaml:"token"`
		ChannelId 	string `yaml:"channel_id"`
		IsUsed		bool `yaml:"is_used"`
	}

	Server struct {
		Port  string `yaml:"port"`
		Mode  string `yaml:"mode"`
		Env   string `yaml:"env"`
		Table string `yaml:"table"`
	} `yaml:"server,omitempty"`
}

var cfg = Config{}
var defaultPath = "./"
var defaultConfigFile = "config/config.yaml"
var ipAddr = ""
var container = ""
var cfgCheckSum = ""

var engine *gin.Engine
var ginLambda *ginadapter.GinLambda

func initServer(launchType string) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rand.Seed(time.Now().UnixNano())
	if launchType != "CONSOLE" {
		awsssm.Init(awsssm.Config{
			Region: "ap-southeast-1",
		})
		cfgCheckSum = awsssm.LoadParamaterStore(
			context.Background(), &cfg, configName)
	} else {
		if len(os.Args) > 1 {
			defaultPath = os.Args[1]
		}
		if len(os.Args) > 2 {
			defaultConfigFile = os.Args[2]
		}

		// unitTest 실행 시 환경설정 참조를 위한 코드
		if len(os.Args) > 3 {
			defaultPath = os.Args[len(os.Args) - 2]
			defaultConfigFile = os.Args[len(os.Args) - 1]
		}

		log.Println(defaultPath + defaultConfigFile)
		core.InitConfig(&cfg, defaultPath+defaultConfigFile)
		cfgCheckSum = core.FileCheckSumOverload(defaultPath + defaultConfigFile)

		// config reload
		go func() {
			c := time.Tick(time.Minute)
			for now := range c {

				checkSum := core.FileCheckSumOverload(defaultPath + defaultConfigFile)
				if cfgCheckSum != checkSum {
					log.Println("config changed: ", now.UTC())
					core.InitConfig(&cfg, defaultPath+defaultConfigFile)
					log.Println(cfg)
					cfgCheckSum = checkSum
				}

				addr := container
				if cfg.Server.Env == "dev" {
					addr = ipAddr
				}
				// server info caching
				redisCache.Set(
					fmt.Sprintf("%s:%s:%s%s", appName, serverName, addr, cfg.Server.Port),
					fmt.Sprintf("port:%s|mode:%s|env:%s", cfg.Server.Port, cfg.Server.Mode, cfg.Server.Env),
					time.Minute*10)
			}
		}()
	}

	log.Println(cfg)
	container, _ = core.GetContainer()
	container = strings.TrimLeft(container, "/")
	if container == "" {
		container = serverName
	}
	ipAddr, _ = core.GetLocalIP()
	log.SetFlags(log.Lshortfile)

	dynamodb.InitDynamoDB(dynamodb.Config{
		Region:    cfg.DynamoDB.Region,
		AccessKey: cfg.DynamoDB.AccessKey,
		SecretKey: cfg.DynamoDB.SecretKey,
		EndPoint:  cfg.DynamoDB.EndPoint,
	})

	awss3.Init(awss3.Config{
		Region:    cfg.S3.Region,
		AccessKey: cfg.S3.AccessKey,
		SecretKey: cfg.S3.SecretKey,
		Bucket:    cfg.S3.Bucket,
	})

	mvlog.InitKinesis(context.Background(), mvlog.KinesisConfig{
		Region:     cfg.Kinesis.Region,
		AccessKey:  cfg.Kinesis.AccessKey,
		SecretKey:  cfg.Kinesis.SecretKey,
		StreamName: cfg.Kinesis.StreamName,
	})

	cacheConfig := redisCache.RedisCacheConfig{
		PAddr:    cfg.Cache.PAddr,
		RAddr:    cfg.Cache.RAddr,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
		Expire:   cfg.Cache.Expire,
		Ver:      cfg.Cache.Ver,
	}

	redisCache.InitializeEx(cacheConfig)

	slack.Init(slack.Config{
		Token:		cfg.Slack.Token,
		ChannelId:	cfg.Slack.ChannelId,
		IsUsed: 	cfg.Slack.IsUsed,
	})

	addr := container
	if cfg.Server.Env == "dev" {
		addr = ipAddr
	}
	// server info caching
	redisCache.Set(
		fmt.Sprintf("%s:%s:%s%s", appName, serverName, addr, cfg.Server.Port),
		fmt.Sprintf("port:%s|mode:%s|env:%s", cfg.Server.Port, cfg.Server.Mode, cfg.Server.Env),
		time.Minute*10)

	router.InitRouter(cfg.Server.Env, defaultPath, ipAddr, cfg.Server.Table, launchType)

	if launchType == "CONSOLE" {
		channel := make(chan os.Signal, 1)
		signal.Notify(channel, syscall.SIGINT)
		signal.Notify(channel, syscall.SIGTERM)

		go func() {
			for sig := range channel {
				fmt.Println("Got signal:", sig)
				addr := container
				if cfg.Server.Env == "dev" {
					addr = ipAddr
				}
				redisCache.Del(
					fmt.Sprintf("%s:%s:%s%s", appName, serverName, addr, cfg.Server.Port))
				time.Sleep(time.Second)
				os.Exit(1)
			}
		}()
	}
}

func routing(lunchType string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(core.Corsm())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	
	if lunchType == "ECS" {
        r.GET("/", func(c *gin.Context) {
            c.JSON(http.StatusOK, nil)
        })
    } else {
        r.GET("/", router.IndexRouter)
    }

	if lunchType != "CONSOLE" {
		r.Use(func(c *gin.Context) {
			// if redisCache.IsConnected() == nil {
			// 	redisCache.SessionKeyToUidHeader(c)
			// } else {
			// 	dynamodb.SessionKeyToUidHeader(c)
			// }

			// uid, err := util.GetHeaderUid(c)
			// if err != nil {
			// 	log.Println(err)
			// 	c.JSON(http.StatusBadRequest, gin.H{
			// 		"err_code": common.ErrCodeShardIndex,
			// 		"err_msg":  err.Error()})
			// 	return
			// }

			// if err = redisCache.BeginProc(uid, c.Request.RequestURI); err != nil {
			// 	log.Println(err)
			// 	c.JSON(http.StatusBadRequest, gin.H{
			// 		"err_code": common.ErrCodeAlreadyProc,
			// 		"err_msg":  err.Error()})
			// 	return
			// }

			if router.TemplateTable == nil {
				router.LoadTable(router.Table)
			}
			
			c.Next()
			
			// redisCache.EndProc(uid)
		})
	}

	// template get/set
	TemplateRouter := r.Group("/template/")
	{
		TemplateRouter.POST("/get", template.GetTemplate)
		TemplateRouter.POST("/set", template.SetTemplate)
	}

	return r
}

// AWS Lambda Handler & Bridge
type HandlerMuti struct{}

func (h *HandlerMuti) Invoke(ctx context.Context, data []byte) ([]byte, error) {
	//log.Println(string(data))

	kEvent := events.APIGatewayProxyRequest{}
	if err := json.Unmarshal(data, &kEvent); err == nil {
		res, err := handlerAPIGateway(ctx, kEvent)
		resData, _ := json.Marshal(res)
		return resData, err
	}

	return handlerEventBridge(ctx, data)
}

func handlerAPIGateway(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func handlerEventBridge(ctx context.Context, data []byte) ([]byte, error) {
	temp := make(map[string]interface{})
	err := json.Unmarshal(data, &temp)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(temp)
	return data, nil
}
// AWS Lambda Handler & Bridge

func main() {
	// launchType AWS ECS | AWS Lambda | Console 에 따라 init
	launchType := "ECS"
	if len(os.Args) > 1 {
		launchType = "CONSOLE"
	}

	if os.Getenv("LAMBDA") == "TRUE" {
		launchType = "LAMBDA"
	}
	log.Println("Launch Type:", launchType)
	initServer(launchType)
	engine = routing(launchType)
	log.Println("Start Template Server port: ", cfg.Server.Port)

	if launchType == "LAMBDA" {
		ginLambda = ginadapter.New(engine)
		lambda.StartHandler(&HandlerMuti{})
	} else {
		engine.Run(cfg.Server.Port)
	}
}
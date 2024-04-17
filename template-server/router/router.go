package router

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	tabletemplate "github.com/chipong/template/common/TableTemplate"

	"github.com/chipong/template/core"
	"github.com/chipong/template/common"
	"github.com/chipong/template/common/util"
	"github.com/chipong/template/common/mvlog"
	"github.com/chipong/template/common/redisCache"
)

var (
	Env 		string
	path		string
	ipaddr		string
	Table		string
	LaunchType	string

	lock		*sync.RWMutex
	serverAddr	map[string]([]string)
	routeServerNames = []string{
		"template-server",
	}

	TemplateTable   *tabletemplate.TemplateData
)

const (
	appName = "template"
	serverName = "template-server"
	s3BucketName = "template"
	tableFolderName = "table/"

	routeUsed = false
)

func InitRouter(serverEnv, defaultPath, ipAddr, table, launchType string) {
	Env = serverEnv
	path = defaultPath
	ipaddr = ipAddr
	Table = table
	LaunchType = launchType

	log.Println("router init -> ", Env, defaultPath, ipAddr)

	serverAddr = make(map[string][]string)

	// 테이블 로딩
	if LaunchType == "ECS" || LaunchType == "CONSOLE" {
		LoadTable(Table)
	}

	if LaunchType != "LAMBDA" {
		go func() {
			c := time.Tick(time.Minute)
			for _ = range c {
				LoadTable(Table)

				// server2server route 사용할 때 목적지 서버 addr set
				if routeUsed {
					lock.Lock()
					for _, serverName := range routeServerNames {
						SetServerAddr(serverName)
					}
					lock.Unlock()
				}
			}
		}()
		// 모니터 틱 시작
		if Env != "dev" {
			go MonitotTick()
		}
	}
}

func LoadTable(table string) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		pc, _, _, _ := runtime.Caller(1)
		log.Println(runtime.FuncForPC(pc).Name()+" took:", elapsed)

	}()

	chs := make([](<-chan *common.ChErrCode), 0)
	ctx, cancel := context.WithCancel(context.Background())

	if table == "S3" {
		loadTemplateTableCh := util.GoRoutineJob(ctx, cancel,
			func(ch chan *common.ChErrCode) {
				defer close(ch)
				if TemplateTable == nil {
					TemplateTable = tabletemplate.NewTemplateData()
				}
				TemplateTable.LoadTableS3(s3BucketName, tableFolderName)
			},
		)
		chs = append(chs, loadTemplateTableCh)
	} else {
		tablePath := strings.Replace(path, `template\`, "", 1)

		loadTemplateTableCh := util.GoRoutineJob(ctx, cancel,
			func(ch chan *common.ChErrCode) {
				defer close(ch)
				if TemplateTable == nil {
					TemplateTable = tabletemplate.NewTemplateData()
				}
				TemplateTable.LoadTable(tablePath + tableFolderName)
			},
		)
		chs = append(chs, loadTemplateTableCh)
	}

	for _, ch := range chs {
		errCh := <-ch
		if errCh != nil {
			panic(errCh.Err.Error())
		}
	}
}

func IndexRouter(c *gin.Context) {
	log.Println("Index Routing")
	LoadTable(Table)
	redisCache.Ping()
	c.JSON(http.StatusOK, nil)
}

func MonitotTick() {
	c := time.Tick(time.Minute)
	for range c {
		monitor := core.MonitorEx()
		mvlog.BeginKinesis()
		mvlog.InsertMVLogMonitor(serverName,
			monitor.Mem, monitor.Cpu, monitor.Disk, monitor.NetSend, monitor.NetRecv)
		mvlog.EndKinesis(context.Background())
	}
}

func SetServerAddr(name string) {
	keys, err := redisCache.Scan(0, appName+fmt.Sprintf(":%s:*", name), 100000)
	if err != nil {
		return
	}

	var urls []string
	for _, v := range keys {
		addr := strings.Replace(v, appName+fmt.Sprintf(":%s:", name), "", 1)
		urls = append(urls, addr)
	}

	serverAddr[name] = urls
}

func GetServerAddr(name string) string {
	lock.RLock()
	defer lock.RUnlock()

	if len(serverAddr[name]) == 0 {
		return ""
	}

	subaddr := ipaddr[:strings.LastIndexByte(ipaddr, '.')]
	addrs := []string{}
	for _, addr := range serverAddr[name] {
		if strings.Contains(addr, subaddr) {
			addrs = append(addrs, addr)
		}
	}

	if len(addrs) == 0 {
		return ""
	}

	return addrs[rand.Intn(len(addrs))]
}
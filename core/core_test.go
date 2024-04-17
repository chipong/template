package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestConfig(t *testing.T) {

	cfg := DefaultConfig{}
	_, err := InitConfig(&cfg, "config/config.yaml")
	if err != nil {
		t.Error(err)
	}
	log.Println(cfg)
}

func TestGinMiddleware(t *testing.T) {
	// setting log(path, name)
	InitLog("", "testcore")

	// new engine
	r := gin.Default()

	// cors
	r.Use(Corsm())

	// logger attribute( skipath, utc )
	r.Use(DefaultLogger())

	// router
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hello Test")
	})

	// run
	go r.Run()

	time.Sleep(time.Second * 5)

	res, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Error(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	log.Println(string(resp))
}

func TestExcel(t *testing.T) {
	sheets, err := ReadExcelSheetMap("cash_gacha.xlsx")
	log.Println("sheets: ", sheets)

	rows, err := ReadExcelFile("cash_gacha.xlsx", "gachaitem")
	if err != nil {
		t.Error(err)
	}

	log.Println("rows: ", len(rows))

	// extract title info
	title := rows[0]
	titleCnt := len(rows[0])

	log.Println(title, titleCnt)
}

func TestGeneratorExcel(t *testing.T) {
	r, err := GeneratorStructuredFile("cash_gacha.xlsx", "gachaitem", "core")
	if err != nil {
		t.Error(err)
	}
	// generator structured go file

	// file delete
	os.Remove("gachaitem.go")

	// save file
	err = ioutil.WriteFile("gachaitem.go", []byte(r), 0)
	if err != nil {
		t.Error(err)
	}
}

func TestDocker(t *testing.T) {
	r, _ := GetContainer()
	log.Println(r)
}

func TestUUID(t *testing.T) {
	//id := guuid.New()
	//log.Println(id.String())

	for i := 0; i < 1; i++ {
		log.Println(GenUUID())
	}

	// for i := 0; i < 1000; i++ {
	// 	log.Println(GenIDNumber())
	// }
}

func TestSessionCheck(t *testing.T) {

	InitConfig(&cfg, "config/config.yaml")
	InitSessionCheck(SessionCheckConfig{
		PAddr:    cfg.Redis.PAddr,
		RAddr:    cfg.Redis.RAddr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	/*
		key, err := GenSessionKey("testaccountid", time.Second*3600*24*365)
		if err != nil {
			log.Println(err)
		}
		log.Println(key)
	*/

	// key, err := ReadRedisClient.Get("f46a637db8e0414bae588250333707f9").Result()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Println(key)

	key, err := GenSessionKey("1234567890", time.Second*3600*24)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(key)

	SessionKeyTTL(key)
}

func TestLocalIP(t *testing.T) {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	// handle err...

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	fmt.Println(localAddr)

	ip, err := GetLocalIP()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ip)
}

func samp(aa interface{}) {
	log.Printf("%v\n", aa)
	newAA := reflect.New(reflect.ValueOf(aa).Elem().Type()).Interface()
	samp2(newAA)
}

func samp2(aa interface{}) {
	log.Printf("%v\n", aa)
}

func TestFileCheckSum(t *testing.T) {
	path := "./../oz/env/dev/table/Account.xlsm"
	checkSum := FileCheckSum(path)

	log.Println("checkSum : ", checkSum)
}

func TestExcelToProto(t *testing.T) {
	fileName := "./../demo/demo-csv/Enums.xlsx"
	outputFileName := "enums.proto"

	err := GenExcelToProtoEnum(fileName, outputFileName)
	if err != nil {
		log.Println(err)
	}
}
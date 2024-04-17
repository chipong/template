package mvlog

import (
	"encoding/json"
	"log"
	"time"
)

const (
	appName = "macovill.oz"
)

type MVLogMonitor struct {
	AppName    string    `json:"app_name"`
	ServerName string    `json:"server_name"`
	CPU        int64     `json:"cpu"`
	Memory     int64     `json:"memory"`
	Disk       int64     `json:"disk"`
	NetSend    int64     `json:"net_send"`
	NetRecv    int64     `json:"net_recv"`
	At         time.Time `json:"at"`
}

func InsertMVLogMonitor(serverName string, mem, cpu, disk, send, recv int64) {
	item := &MVLogMonitor{
		AppName:    appName,
		ServerName: serverName,
		CPU:        cpu,
		Memory:     mem,
		Disk:       disk,
		NetSend:    send,
		NetRecv:    recv,
		At:         time.Now().UTC(),
	}
	//core.InsertMVLogItem(item)
	jsonData, _ := json.Marshal(item)
	//log.Println(string(jsonData))
	WriteBatch("MVLogMonitor", jsonData)
	PutKinesisMulti("MVLogMonitor", jsonData)
}

// MVLogPing ...
type MVLogPing struct {
	AppName string      `json:"app_name"`
	UID     string      `json:"uid"`
	Screen  interface{} `json:"screen"`
	StartAt int64       `json:"start_at"`
	At      time.Time   `json:"at"`
}

// InsertMVLogPing ...
func InsertMVLogPing(appName, uid string, screen interface{}, startat int64) {
	item := &MVLogPing{
		AppName: appName,
		UID:     uid,
		Screen:  screen,
		StartAt: startat,
		At:      time.Now().UTC(),
	}
	//core.InsertMVLogItem(item)
	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogPing", jsonData)
	PutKinesisMulti("MVLogPing", jsonData)
}

type MVLogPost struct {
	AppName string `json:"app_name"` // 앱 이름
	Uid     string `json:"uid"`
	Cn      string `json:"cn"` // 컨텐츠 네임

	Pid      int64     `json:"pid"`       // 우편 고유 번호
	PostType int64     `json:"post_type"` // 우편 타입
	Status   int64     `json:"status"`    // 상태
	Coin     int64     `json:"coin"`      // 첨부 코인
	CashFree int64     `json:"cash_free"` // 첨부 캐시
	CashIap  int64     `json:"cash_iap"`  // 첨부 캐시
	CardId   int64     `json:"card_id"`   // 첨부 카드
	Count    int64     `json:"count"`     // 첨부 카드 갯수
	At       time.Time `json:"at"`
}

func InsertMVLogPost(uid, cn string, pid, posttype, status, coin, cashfree, cashiap, cardid, count int64) {
	item := &MVLogPost{
		AppName:  appName,
		Uid:      uid,
		Cn:       cn,
		Pid:      pid,
		PostType: posttype,
		Status:   status,
		Coin:     coin,
		CashFree: cashfree,
		CashIap:  cashiap,
		CardId:   cardid,
		Count:    count,
		At:       time.Now().UTC(),
	}
	//core.InsertMVLogItem(item)
	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogPost", jsonData)
	PutKinesisMulti("MVLogPost", jsonData)
}

type MVLogIap struct {
	AppName   string    `json:"app_name"` // 앱 이름
	Uid       string    `json:"uid"`
	Cn        string    `json:"cn"`         // 컨텐츠 네임
	ProductId string    `json:"product_id"` // 상품 아이디
	Rid       int64     `json:"rid"`        // 영수증 고유 번호
	Status    int64     `json:"status"`     // 상태
	PayType   int64     `json:"pay_type"`   // 구매 타입
	At        time.Time `json:"at"`
}

func InsertMVLogIap(uid, cn, productid string, rid, status, paytype int64) {
	item := &MVLogIap{
		AppName:   appName,
		Uid:       uid,
		Cn:        cn,
		ProductId: productid,
		Rid:       rid,
		Status:    status,
		PayType:   paytype,
		At:        time.Now().UTC(),
	}
	//core.InsertMVLogItem(item)
	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogIap", jsonData)
	PutKinesisMulti("MVLogIap", jsonData)
}

type MVLogSnap struct {
	AppName string `json:"app_name"` // 앱 이름
	Uid     string `json:"uid"`
	Cn      string `json:"cn"` // 컨텐츠 네임

	Account     interface{} `json:"account"`      // 계정 정보
	PointAsstes interface{} `json:"point_assets"` // 포인트 제화 정보
	Char        interface{} `json:"char"`         // 캐릭터 정보(현재 댁의 카드만)
	Stage       interface{} `json:"stage"`        // 스테이지 클리어 정보

	At time.Time `json:"at"`
}

func InsertMVLogSnap(uid, cn string, account, pointassets, char, stage interface{}) {
	item := &MVLogSnap{
		AppName:     appName,
		Uid:         uid,
		Cn:          cn,
		Account:     account,
		PointAsstes: pointassets,
		Char:        char,
		Stage:       stage,
		At:          time.Now().UTC(),
	}

	//core.InsertMVLogItem(item)
	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogSnap", jsonData)
	PutKinesisMulti("MVLogSnap", jsonData)
}

type MVLogStage struct {
	AppName string `json:"app_name"` // 앱 이름
	Uid     string `json:"uid"`
	Cn      string `json:"cn"` // 컨텐츠 네임

	StageId      string    `json:"stage_id"`
	ClearStatus  int64     `json:"clear_status"`
	RewardStatus int64     `json:"reward_status"`
	ClearSec     int64     `json:"clear_sec"`
	PartyNo      int64     `json:"party_no"`
	At           time.Time `json:"at"`
}

func InsertMVLogStage(uid, cn string, stageid string, clear, reward, sec, partyno int64) {
	item := &MVLogStage{
		AppName:      appName,
		Uid:          uid,
		Cn:           cn,
		StageId:      stageid,
		ClearStatus:  clear,
		RewardStatus: reward,
		ClearSec:     sec,
		PartyNo:      partyno,
		At:           time.Now().UTC(),
	}
	//core.InsertMVLogItem(item)
	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogStage", jsonData)
	PutKinesisMulti("MVLogStage", jsonData)
}

type MVLogInvenStackable struct {
	AppName string `json:"app_name"` // 앱 이름
	Uid     string `json:"uid"`
	Cn      string `json:"cn"` // 컨텐츠 네임

	ItemId     string    `json:"item_id,omitempty"`     // 아이템 아이디
	Count      int64     `json:"count,omitempty"`       // 아이템 갯수 증감
	AfterCount int64     `json:"after_count,omitempty"` // 최종 아이템 갯수
	At         time.Time `json:"at"`
}

func InsertMVLogInvenStackable(uid, cn string, itemid string, count, aftercount int64) {
	item := &MVLogInvenStackable{
		AppName:    appName,
		Uid:        uid,
		Cn:         cn,
		ItemId:     itemid,
		Count:      count,
		AfterCount: aftercount,
		At:         time.Now().UTC(),
	}
	//core.InsertMVLogItem(item)
	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogInvenStackable", jsonData)
	PutKinesisMulti("MVLogInvenStackable", jsonData)
}

type MVLogConsumeCash struct {
	AppName string `json:"app_name"` // 앱 이름
	Uid     string `json:"uid"`
	Cn      string `json:"cn"` // 컨텐츠 네임

	TargetContents string `json:"target_contents,omitempty"` // 대상 컨텐츠

	Cash     int64 `json:"cash,omitempty"`      // 소모 캐쉬
	CashFree int64 `json:"cash_free,omitempty"` // 소모 무료 캐쉬

	AfterCash     int64     `json:"after_cash,omitempty"`      // 최종 캐쉬
	AfterCashFree int64     `json:"after_cash_free,omitempty"` // 최종 무료 캐쉬
	At            time.Time `json:"at"`
}

func InsertMVLogConsumeCash(uid, cn string, contents string, cash, cashfree, aftercash, aftercashfree int64) {
	item := &MVLogConsumeCash{
		AppName:        appName,
		Uid:            uid,
		Cn:             cn,
		TargetContents: contents,
		Cash:           cash,
		CashFree:       cashfree,
		AfterCash:      aftercash,
		AfterCashFree:  aftercashfree,
		At:             time.Now().UTC(),
	}
	//core.InsertMVLogItem(item)
	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogConsumeCash", jsonData)
	PutKinesisMulti("MVLogConsumeCash", jsonData)
}

type MVLogSummonEx struct {
	AppName		string `json:"app_name"` // 앱 이름
	Uid			string `json:"uid"`
	Cn			string `json:"cn"` // 컨텐츠 네임

	SummonId   string    `json:"summon_id,omitempty"`	// 소환 아이디
	Count      int64     `json:"count,omitempty"`       // 소환 횟수
	AfterCount int64     `json:"after_count,omitempty"`	// 최종 소환 횟수
	At         time.Time `json:"at"`
}

func InsertMVLogSummonEx(uid, cn string, summonid string, count, aftercount int64) {
	item := &MVLogSummonEx{
		AppName:    appName,
		Uid:        uid,
		Cn:         cn,
		SummonId:   summonid,
		Count:      count,
		AfterCount: aftercount,
		At:         time.Now().UTC(),
	}
	
	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogSummon", jsonData)
	PutKinesisMulti("MVLogSummon", jsonData)
}

type MVLogNewAccount struct {
	AppName string `json:"app_name"` // 앱 이름
	Uid     string `json:"uid"`
	Cn      string `json:"cn"` // 컨텐츠 네임

	Nickname string    `json:"nickname"` // 닉네임
	At       time.Time `json:"at"`
}

func InsertMVLogNewAccountEx(
	uid, cn string,
	nickname string) {

	item := &MVLogNewAccount{
		AppName:  appName,
		Uid:      uid,
		Cn:       cn,
		Nickname: nickname,

		At: time.Now().UTC(),
	}

	jsonData, _ := json.Marshal(item)
	log.Println(string(jsonData))
	WriteBatch("MVLogNewAccount", jsonData)
	PutKinesisMulti("MVLogNewAccount", jsonData)
}
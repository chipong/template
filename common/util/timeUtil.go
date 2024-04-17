package util

import (
	"log"
	"time"
	"strings"
	"strconv"
)

/*
const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)
*/

const (
	dateLayout = "2006-01-02 15:04:05"
)

var BaseTimePeriod int64 = 0 // utc +0

func OzNow() time.Time {
	return time.Now().Add(time.Duration(BaseTimePeriod * int64(time.Second)))
}

func OzNowUnix() int64 {
	return OzNow().Unix()
}

func SetBaseTimePeriod(period int64) {
	BaseTimePeriod = period
	log.Println("Base Time Add: ", period, OzNow())
}

func TodayMidnightTime() time.Time {
	return time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(),
		time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC)
}

func TodayByHMS(hh_mm_ss string) time.Time {
	splitHMS := strings.Split(hh_mm_ss, ":")
	hour, _ := strconv.Atoi(splitHMS[0])
	minete, _ := strconv.Atoi(splitHMS[1])
	second, _ := strconv.Atoi(splitHMS[2])
	return time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(),
	time.Now().UTC().Day(), hour, minete, second, 0, time.UTC)
}

func TomorrowMidnightTime() time.Time {
	now := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(),
		time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC)
	return now.AddDate(0, 0, 1)
}

func YesterdayMidnightTime() time.Time {
	now := time.Date(time.Now().UTC().Year(), time.Now().UTC().Month(),
		time.Now().UTC().Day(), 0, 0, 0, 0, time.UTC)
	return now.AddDate(0, 0, -1)
}

func DayChange(timestamp int64) bool {
	return time.Unix(timestamp, 0).UTC().Day() != time.Now().UTC().Day()
}

// Golang week 시작은 MON
func WeekChange(timestamp int64) bool {
	_, w := time.Unix(timestamp, 0).UTC().ISOWeek()
	_, nowWeek := time.Now().UTC().ISOWeek()
	return w != nowWeek
}

// func DayChangeWithResetTime(timestamp int64, resettime int64) bool {
// 	timestamp = timestamp - resettime*3600
// 	return time.Unix(timestamp, 0).UTC().Day() != time.Now().UTC().Day()
// }

func DayChangeWithResetTime(timestamp int64, resettime string) bool {
	resetTimestamp := TodayByHMS(resettime)
	log.Println("timestamp : ", time.Unix(timestamp, 0).Unix())
	log.Println("resetTimestamp : ", resetTimestamp.Unix())
	return time.Unix(timestamp, 0).Unix() < resetTimestamp.Unix()
}

func StringToZoneTimeStamp(strDate string) (int64, error) {
	date, err := time.Parse(dateLayout, strDate)
	if err != nil {
		return 0, err
	}

	t := time.Now()
    _, offset := t.Zone()

	timeStamp := date.Unix() - int64(offset)
	
	return timeStamp, nil
}

func StringToTimeStamp(strDate string) (int64, error) {
	date, err := time.Parse(dateLayout, strDate)
	if err != nil {
		return 0, err
	}
	
	return date.Unix(), nil
}

// 이벤트 시작-끝 시간을 통해 사용여부 체크
func CheckAvailableDate(start, end string) (bool, error) {
	startDate, err := time.Parse(dateLayout, start)
	if err != nil {
		return false, err
	}

	endDate, err := time.Parse(dateLayout, end)
	if err != nil {
		return false, err
	}

	t := time.Now()
    // _, offset := t.Zone()

	nowTimestamp := t.Unix()
	startTimestamp := startDate.Unix()
	endTimestamp := endDate.Unix()
	// startTimestamp := startDate.Unix() - int64(offset)
	// endTimestamp := endDate.Unix() - int64(offset)

	// log.Printf("nowTimestamp	: %d\n", nowTimestamp)
	// log.Printf("startTimestamp	: %d\n", startTimestamp)
	// log.Printf("endTimestamp	: %d\n", endTimestamp)

	// log.Printf("nowTimestamp < startTimestamp	: %v", nowTimestamp < startTimestamp)
	// log.Printf("nowTimestamp >= endTimestamp	: %v", nowTimestamp >= endTimestamp)

	if nowTimestamp < startTimestamp || nowTimestamp >= endTimestamp {
		return false, nil
	}
    
	return true, nil
}
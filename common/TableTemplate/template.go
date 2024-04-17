package tabletemplate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/chipong/template/core"
	awss3 "github.com/chipong/template/common/awsS3"
	"github.com/chipong/template/common/proto"
	"github.com/chipong/template/common/slack"
)

const (
	HEADER_ROW_COUNT = 6
	START_COLUMN     = 1
)

type TemplateData struct {
	temp			map[string]*oz.Template
	checkSumTemp	string
	lock			*sync.RWMutex
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		temp:			make(map[string]*oz.Template),
		checkSumTemp:	"",
		lock:			new(sync.RWMutex),
	}
}

func (t *TemplateData) LoadTable(path string) error {
	{
		checkSum := core.FileCheckSumOverload(path + "Template.xlsm")
		if t.checkSumTemp != "" && t.checkSumTemp == checkSum {
			return nil
		}

		{
			sheetMap, err := core.ReadExcelSheetMap(path + "Template.xlsm")
			if err != nil {
				log.Println(err.Error())
				return err
			}
	
			sheet := "Template"
			table := [][]string{}
			for _, sheetName := range sheetMap {
				if strings.Compare(sheetName, "Info") == 0 || !strings.Contains(sheetName, sheet) {
					continue
				}
	
				tableTemp, err := core.ReadExcelFile(path+"Template.xlsm", sheetName)
				if err != nil {
					log.Println(err.Error())
					return err
				}
				
				var isEmpty bool = false
				if len(table) == 0 {
					isEmpty = true
				}
	
				// 현재 table 목록이 비어있으면 header append
				for index, v := range tableTemp {
					if index < HEADER_ROW_COUNT && !isEmpty {
						continue
					}
	
					table = append(table, v)
				}
			}
	
			core.UpperTable(table)
			t.loadTemplate(table, sheet)
		}
	
	
		t.checkSumTemp = checkSum
	}

	return nil
}

func (t *TemplateData) LoadTableS3(bucket, path string) error {
	{
		data, err := awss3.GetObject(context.Background(), bucket, path+"Template.xlsm")
		if err != nil {
			log.Println(err)
			return err
		}
	
		checkSum := core.FileCheckSumOverload(data)
		if t.checkSumTemp != "" && t.checkSumTemp == checkSum {
			return nil
		}
	
		log.Println("aws ", bucket, path)
		{
			closeReader := io.NopCloser(bytes.NewReader(data))
			defer closeReader.Close()
	
			sheetMap, err := core.ReadExcelReaderSheetMap(closeReader)
			if err != nil {
				log.Println(err.Error())
				return err
			}
	
			sheet := "Template"
			table := [][]string{}
			for _, sheetName := range sheetMap {
				if strings.Compare(sheetName, "Info") == 0 || !strings.Contains(sheetName, sheet) {
					continue
				}
	
				cReader := io.NopCloser(bytes.NewReader(data))
				defer cReader.Close()
	
				tableTemp, err := core.ReadExcelReader(cReader, sheetName)
				if err != nil {
					log.Println(err.Error())
					return err
				}
	
				var isEmpty bool = false
				if len(table) == 0 {
					isEmpty = true
				}
	
				// 현재 table 목록이 비어있으면 header append
				for index, v := range tableTemp {
					if index < HEADER_ROW_COUNT && !isEmpty {
						continue
					}
	
					table = append(table, v)
				}
			}
	
			core.UpperTable(table)
			t.loadTemplate(table, sheet)
		}
	
		t.checkSumTemp = checkSum
	}

	return nil
}

func (t *TemplateData) loadTemplate(table [][]string, sheet string) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	tempTable := make(map[string]*oz.Template)
	for index, col := range table {
		// header skip
		if index < HEADER_ROW_COUNT {
			continue
		}
		if strings.Contains(col[1], "//") {
			log.Println("comment-> ", col)
			continue
		}
		i := 1

		id := col[Inc(&i)]
		groupId := col[Inc(&i)]
		rewardId := col[Inc(&i)]

		count, err := core.ParseInt(col[Inc(&i)])
		if err != nil {
			TableLoadErr(sheet, table[4][i - 1], index, i, err)
			return err
		}

		templateEnum := oz.TemplateEnum_T(oz.TemplateEnum_T_value[col[Inc(&i)]])
		if templateEnum == oz.TemplateEnum_NONE || templateEnum == oz.TemplateEnum_MAX {
			log.Println(col)
			err = errors.New("enum range over")
			TableLoadErr(sheet, table[4][i - 1], index, i, err)
			return err
		}

		tempTable[id] = &oz.Template{
			Id:			id,
			GroupId: 	groupId,
			RewardId:	rewardId,
			Count:		int64(count),
			Enum:		templateEnum,
		}
	}

	t.temp = tempTable
	log.Println("table loaded ", sheet, len(t.temp))
	return nil
}

func (t *TemplateData) FindTemplate(id string) (*oz.Template, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	id = strings.ToUpper(id)
	if t.temp[id] == nil {
		return nil, errors.New("not found data")
	}
	return t.temp[id], nil
}

func Inc(index *int) int {
	i := *index
	*index++
	return i
}

func TableLoadErr(table, columnName string, row, col int, err error) {
	errMsg := fmt.Sprintf("%s %s row:%d col:%d err:%s", table, columnName, row, col, err.Error())
	log.Println(errMsg)
	if slack.GetConfig().IsUsed {
		slack.SendMessage(errMsg)
	}
	panic(err)
}
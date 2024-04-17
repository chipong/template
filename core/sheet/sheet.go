package sheet

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

var (
	service *spreadsheet.Service
)

func InitSpreadSheet(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return
	}

	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		log.Println("jwt", err)
		return
	}

	client := conf.Client(context.TODO())
	service = spreadsheet.NewServiceWithClient(client)
}

func FetchSpreadSheet(sheetIdentity string, sheetIndex uint) ([][]string, error) {
	if service == nil {
		return nil, errors.New("service is not connected")
	}

	ssheet, err := service.FetchSpreadsheet(sheetIdentity)
	if err != nil {
		log.Println("fetch", err)
		return nil, err
	}

	sheet, err := ssheet.SheetByIndex(sheetIndex)
	if err != nil {
		log.Println("sheet", err)
		return nil, err
	}

	datas := make([]([]string), 0)
	for _, row := range sheet.Rows {
		rowCell := make([]string, 0)
		for _, cell := range row {
			value := fmt.Sprintf("\"%s\"", cell.Value)
			rowCell = append(rowCell, value)
		}
		datas = append(datas, rowCell)
	}
	return datas, nil
}

func FetchSpreadSheetById(sheetIdentity string, sheetTitle string) ([][]string, error) {
	if service == nil {
		return nil, errors.New("service is not connected")
	}

	ssheet, err := service.FetchSpreadsheet(sheetIdentity)
	if err != nil {
		log.Println("fetch", err)
		return nil, err
	}

	sheet, err := ssheet.SheetByTitle(sheetTitle)
	if err != nil {
		log.Println("sheet", err)
		return nil, err
	}

	datas := make([]([]string), 0)
	for _, row := range sheet.Rows {
		rowCell := make([]string, 0)
		for _, cell := range row {
			value := fmt.Sprintf("\"%s\"", cell.Value)
			rowCell = append(rowCell, value)
		}
		datas = append(datas, rowCell)
	}
	return datas, nil
}

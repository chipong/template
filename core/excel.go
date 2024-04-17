package core

import (
	"os"
	"io"
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// ReadExcelFile ...
func ReadExcelFile(file string, sheet string) ([][]string, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, err
	}

	rows := f.GetRows(sheet)
	// for _, row := range rows {
	// 	for _, col := range row {
	// 		if col == "" {
	// 			continue
	// 		}
	// 		log.Println(col)
	// 	}
	// 	log.Println()
	// }
	return rows, nil
}

func ReadExcelReader(r io.ReadCloser, sheet string) ([][]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	rows := f.GetRows(sheet)
	// for _, row := range rows {
	// 	for _, col := range row {
	// 		if col == "" {
	// 			continue
	// 		}
	// 		log.Println(col)
	// 	}
	// 	log.Println()
	// }
	return rows, nil
}

// ReadExcelSheetMap ...
func ReadExcelSheetMap(file string) (map[int]string, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, err
	}

	return f.GetSheetMap(), nil
}

// ReadExcelReaderSheetMap ...
func ReadExcelReaderSheetMap(r io.ReadCloser) (map[int]string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}

	return f.GetSheetMap(), nil
}

// GeneratorStructuredFile ...
func GeneratorStructuredFile(file string, sheet string, pkg string) (string, error) {
	rows, err := ReadExcelFile(file, sheet)
	if err != nil {
		return "", err
	}

	structureName := strings.ToUpper(sheet)

	// extract title info
	title := rows[0]

	resultStr := "package " + pkg + "\r\n"
	resultStr += "// " + structureName + " ... \r\n"
	resultStr += "type " + structureName + " struct {\r\n"
	for _, row := range title {
		resultStr += "\t" + strings.ToUpper(row) + "\tstring `json:\"" + row + "\"`\r\n"
	}
	resultStr += "}\r\n"

	resultStr += "// Load" + structureName + " ... \r\n"
	resultStr += "func Load" + structureName + "(rows [][]string) interface{} {\r\n"
	resultStr += "\titems := make([]" + structureName + ", 0)\r\n\r\n"

	resultStr += "\tfor _, row := range rows[1:] {\r\n"
	resultStr += "\t\titem := " + structureName + "{}\r\n"

	for i, row := range title {
		resultStr += "\t\titem." + strings.ToUpper(row) + " = row[" + strconv.Itoa(i) + "]\r\n"
	}

	resultStr += "\t\titems = append(items, item)\r\n"
	resultStr += "\t}\r\n"
	resultStr += "\treturn items\r\n"
	resultStr += "}\r\n"

	return resultStr, nil
}

func GenExcelToProtoEnum(fileName, outputFileName string) error {
    proto := "syntax = \"proto3\";\n\noption go_package = \"./oz\";\noption java_multiple_files = true;\n\npackage oz;\n\n"

    sheetMap, err := ReadExcelSheetMap(fileName)
    if err != nil {
        return err
    }
    
    for _, sheetName := range sheetMap {
        enums, err := ReadExcelFile(fileName, sheetName)
        if err != nil {
            return err
        }
    
        lenRows := len(enums)
        lenColumns := len(enums[0])
        for c := 0; c < lenColumns; c++ {
            emptyCount := 0
            for r := 0; r < lenRows; r++ {
                value := enums[r][c];
                if value == "" {
                    emptyCount++
                    continue
                }
    
                if r == 0 {
                    proto += fmt.Sprintf("message %s {\n\tenum T {\n\t\tNONE = 0;\n", value)
                } else {
                    proto += fmt.Sprintf("\t\t%s = %d;\n", value, r)
                }
            }
    
            proto += fmt.Sprintf("\t\tMAX = %d;\n\t}\n}\n\n", lenRows - emptyCount)
        }
    }

    f, err := os.Create(outputFileName)
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = f.WriteString(proto)
    if err != nil {
        return err
    }

    f.Sync()
    
    return nil
}
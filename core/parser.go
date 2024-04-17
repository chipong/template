package core

import (
	"strconv"
	"strings"
	"time"
)

func UpperTable(table [][]string) [][]string {
	for index, col := range table {
		for r, row := range col {
			row = strings.ReplaceAll(row, "\"", "")
			if row == "-" {
				row = ""
			}
			table[index][r] = strings.ToUpper(row)
		}
	}
	return table
}

func ParseFloat(str string) (float64, error) {
	if str == "" {
		return 0, nil
	}
	return strconv.ParseFloat(str, 64)
}

func ParseInt64(str string) (int64, error) {
	if str == "" {
		return 0, nil
	}
	return strconv.ParseInt(str, 10, 64)
}

func ParseInt(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	return strconv.Atoi(str)
}

func ParseBool(str string) (bool, error) {
	if str == "" {
		return false, nil
	}
	return strconv.ParseBool(str)
}

func ParseDateTime(str string) (int64, error) {
	if str == "" {
		return 0, nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func Inc(index *int) int {
	*index++
	return *index
}

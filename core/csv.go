package core

import (
	"os"
	"bytes"
	"encoding/csv"
)

func ReadCsvFile(fileName string) ([][]string, error) {
    f, err := os.Open(fileName)
    if err != nil {
        return [][]string{}, err
    }
    defer f.Close()

    r := csv.NewReader(f)
	r.Comma = ','
	r.LazyQuotes = true

    rows, err := r.ReadAll()

    if err != nil {
        return [][]string{}, err
    }

    return rows, nil
}

func ReadCsvReader(fileName string, body []byte) ([][]string, error) {
	r := csv.NewReader(bytes.NewBuffer(body))
	r.Comma = ','
	r.LazyQuotes = true

    rows, err := r.ReadAll()

    if err != nil {
        return [][]string{}, err
    }

    return rows, nil
}
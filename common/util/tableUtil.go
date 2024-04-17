package util

import (
	// "errors"
	"io/fs"
	"io/ioutil"
	"log"
)

func GetFiles(path string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return files, nil
}
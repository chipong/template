package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// DefaultConfig ...
type DefaultConfig struct {
	DynamoDB struct {
		Region    string `yaml:"region"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_access_key"`
		EndPoint  string `yaml:"end_point"`
	} `yaml:"dynamodb,omitempty"`

	Mysql struct {
		Addr string `yaml:"addr"`
		User string `yaml:"user"`
		PWD  string `yaml:"pwd"`
		DB   string `yaml:"db"`
	} `yaml:"mysql,omitempty"`

	Redis struct {
		PAddr    string `yaml:"primary_addr"`
		RAddr    string `yaml:"reader_addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis,omitempty"`

	Server struct {
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server,omitempty"`
}

// cfg ...
var cfg = DefaultConfig{}

// InitConfig ...
func InitConfig(conf interface{}, path string) (interface{}, error) {
	filename, _ := filepath.Abs(path)
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = yaml.Unmarshal(f, conf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return conf, nil
}

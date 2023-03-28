package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type configFormat struct {
	Address       string `yaml:"Address"`
	MysqlAddress  string `yaml:"Mysql_address"`
	MysqlUsername string `yaml:"Mysql_username"`
	MysqlPassword string `yaml:"Mysql_password"`
	MysqlDatabase string `yaml:"Mysql_database"`
}

var config configFormat

func loadConfig(configFile string) configFormat {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err)
	}
	if err := config.Parse(data); err != nil {
		log.Println(err)
	}
	return config
}

func (c *configFormat) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

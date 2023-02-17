package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type configFormat struct {
	Address        string `yaml:"Address"`
	Mysql_address  string `yaml:"Mysql_address"`
	Mysql_username string `yaml:"Mysql_username"`
	Mysql_password string `yaml:"Mysql_password"`
	Mysql_database string `yaml:"Mysql_database"`
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

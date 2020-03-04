package ue_factory

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func checkErr(err error) {
	if err != nil {
		err = fmt.Errorf("[Configuration] %s", err.Error())
		log.Panic(err.Error())
	}
}

func InitUeConfigFactory(f string) *Config {
	content, err := ioutil.ReadFile(f)
	checkErr(err)

	config := Config{}

	err = yaml.Unmarshal([]byte(content), &config)
	checkErr(err)
	return &config
}

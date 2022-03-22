package config

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type YAML struct {
	Port          string `yaml:"port"`
	Host          string `yaml:"host"`
	Mailsender    string `yaml:"mailsender"`
	Redisdb       int    `yaml:"redisdb"`
	Redisport     string `yaml:"redisport"`
	Redispassword string `yaml:"redispassword"`
	Imagespath    string `yaml:"imagespath"`
}

var C YAML

func Init() {
	yamlFile, _ := ioutil.ReadFile("pkg/config/conf.yaml")
	err := yaml.Unmarshal(yamlFile, &C)
	if err != nil {
		logrus.Error(err)
	}
}

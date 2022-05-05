package configuration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Maps struct {
	Provider string `yaml:"provider"`
	ApiKey   string `yaml:"api_key"`
}

type App struct {
	Port  int  `yaml:"port"`
	Debug bool `yaml:"debug"`
}

type Configuration struct {
	MAPS Maps `yaml:"maps"`
	APP  App  `yaml:"app"`
}

const defaultPath string = "config.yaml"

func Load(filename string) (*Configuration, error) {
	config := &Configuration{}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func New() (*Configuration, error) {
	return Load(defaultPath)
}

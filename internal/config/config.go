package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Directory  string   `json:"directory"`
	Host       string   `json:"host"`
	Port       string   `json:"port"`
	Extensions []string `json:"extensions"`
}

func NewConfig(filename string) *Config {
	c := &Config{}
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(text, &c)
	if err != nil {
		panic(err)
	}
	return c
}

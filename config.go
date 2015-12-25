package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Server        string `json:"server"`
	Port          int    `json:"port"`
	RedisAddr     string `json:"redis_addr"`
	RedisPort     int    `json:"redis_port"`
	RedisDb       int    `json:"redis_db"`
	RedisPassword string `json:"redis_password"`
}

func ParseConfig(path string) (config *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	config = &Config{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return
}

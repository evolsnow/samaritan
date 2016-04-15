package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type listenServer struct {
	Address  string `json:"address,omitempty"`
	Port     int    `json:"port"`
	DB       string `json:"db,omitempty"`
	Password string `json:"password,omitempty"`
}

// Config is a web server config
type Config struct {
	WebS   listenServer `json:"web_server"`
	RedisS listenServer `json:"redis_server"`
	MysqlS listenServer `json:"mysql_server"`
	RpcSD  listenServer `json:"rpc_server_d"`
	RpcSF  listenServer `json:"rpc_server_f"`
}

// ParseConfig parses config from the given file path
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

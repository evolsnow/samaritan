package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ListenServer struct {
	Address  string `json:"address,omitempty"`
	Port     int    `json:"port"`
	DB       string `json:"db,omitempty"`
	Password string `json:"password,omitempty"`
}

type Config struct {
	HttpS  ListenServer `json:"http_server"`
	RedisS ListenServer `json:"redis_server"`
	MysqlS ListenServer `json:"mysql_server"`
	RpcSD  ListenServer `json:"rpc_server_d"`
	RpcSF  ListenServer `json:"rpc_server_f"`
}

// ParseConfig parse config from the given file path
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

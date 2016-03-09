package main

import "testing"

func TestParseConfig(t *testing.T) {
	cfg, err := ParseConfig("config.json.example")
	if err != nil {
		t.Error("parse config.json errer:", err.Error())
	}
	if cfg.MysqlPassword == "" {
		t.Error("parse struct error")
	}
	_, err = ParseConfig("non-existen.json")
	if err == nil {
		t.Error("should throw err when parse from non-existen json file")
	}
}

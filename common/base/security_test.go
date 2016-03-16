package base

import (
	"github.com/evolsnow/httprouter"
	"testing"
)

var ah string
var raw = "raw"

func TestNewToken(t *testing.T) {
	ah = MakeToken(123)
	if ah == "" {
		t.Error("generate new token error")
	}
}

func TestParseToken(t *testing.T) {
	ps := new(httprouter.Params)
	if err := ParseToken(ah, ps); err != nil {
		t.Error("parse token error:", err)
	}
	if ps.GetInt("authId") != 123 {
		t.Error("token parse result mismatch")
	}
	if err := ParseToken("invalid token", ps); err == nil {
		t.Error("parse invalid token error")
	}
}

func TestEncryptPassword(t *testing.T) {
	enc := EncryptedPassword(raw)
	if enc == "" || enc == raw {
		t.Error("encrypt password error")
	}
}

func TestHashWithSalt(t *testing.T) {
	enc := hashWithSalt(raw, "salt")
	if enc == "" || enc == raw {
		t.Error("hash with salt error")
	}
}

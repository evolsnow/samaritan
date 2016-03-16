package base

import (
	"testing"
)

var ah string
var raw = "raw"

func TestNewToken(t *testing.T) {
	ah = MakeToken(123)
	if ah == "" {
		t.Error("generate new token error")
	}
	ah = ah[7:]
}

func TestParseToken(t *testing.T) {
	var uid int
	var err error
	if uid, err = ParseToken(ah); err != nil {
		t.Error("parse token error:", err)
	}
	if uid != 123 {
		t.Error("token parse result mismatch")
	}
	if _, err := ParseToken("invalid token"); err == nil {
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

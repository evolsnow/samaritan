package base

import (
	"strconv"
	"testing"
)

func TestRandomCode(t *testing.T) {
	code := RandomCodeSix()
	if len(code) != 6 {
		t.Error("code length mismatch")
	}
	if _, err := strconv.Atoi(code); err != nil {
		t.Error("code is not raw number")
	}
}

func TestValidPhone(t *testing.T) {
	if !ValidPhone("13212345678") {
		t.FailNow()
	}
	if ValidPhone("1329123478") {
		t.FailNow()
	}
	if ValidPhone("11011011011") {
		t.FailNow()
	}
}

func TestValidMail(t *testing.T) {
	if !ValidMail("foo@bar.com") {
		t.FailNow()
	}
	if ValidMail("foo.bar") {
		t.FailNow()
	}
	if ValidMail("foo@") {
		t.FailNow()
	}
	if ValidMail("@bar.com") {
		t.FailNow()
	}
}

func TestValidSamId(t *testing.T) {
	if !ValidSamId("valid_id") {
		t.Error("valid_id")
	}
	if !ValidSamId("0111") {
		t.Error("0111 valid")
	}
	if !ValidSamId("abc") {
		t.Error("abc valid")
	}
	if ValidSamId("invalid#!id") {
		t.Error("middle")
	}
	if ValidSamId("*%invalid#!id") {
		t.Error("begin")
	}
	if ValidSamId("") {
		t.Error("empty")
	}
}

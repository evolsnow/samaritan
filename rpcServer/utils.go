package main

import (
	"fmt"
	"math/rand"
	"net/mail"
	"strings"
	"time"
)

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), "<@>")
}

func makeMessageId(domain string) string {
	now := time.Now()
	utcDate := now.Format("20060102150405")
	rdm := rand.New(rand.NewSource(now.UnixNano()))
	randInt := rdm.Intn(100000)
	return fmt.Sprintf("<%d.%d@%s>", utcDate, randInt, domain)
}

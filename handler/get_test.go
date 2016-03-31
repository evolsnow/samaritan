package handler

import (
	"encoding/json"
	"github.com/evolsnow/samaritan/common/dbms"
	"github.com/evolsnow/samaritan/common/log"
	"net/http"
	"testing"
)

func init() {
	dbms.Pool = dbms.NewPool("127.0.0.1:6379", "", "1")
	//c := dbms.Pool.Get()
	//defer c.Close()
	//c.Do("FLUSHDB")
}

func get(reqURL string, ds interface{}) {
	var t testing.T
	//reqURL = url.QueryEscape(reqURL)
	log.Info(reqURL)
	resp, err := http.Get(reqURL)
	if err != nil {
		t.Error("http get err")
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(ds)
	if err != nil {
		t.Error(err)
	}
}

func TestSamIdStatus(t *testing.T) {

	reply := new(samIdStatusResp)

	dbms.DeleteSamId("testevol")
	get("http://127.0.0.1:8080/samIds/testevol", reply)
	if !reply.Available {
		t.Error("should be available")
	}

	dbms.UpdateSamIdSet("testevol")
	get("http://127.0.0.1:8080/samIds/testevol", reply)
	if reply.Available || reply.Msg != ExistErr {
		t.Error("should be unavailable")
	}

	get(`http://127.0.0.1:8080/samIds/*!1234`, reply)
	if reply.Msg != CharsetErr {
		t.Error("illegal charset")
	}

	get("http://127.0.0.1:8080/samIds/abc", reply)
	if reply.Msg != LengthErr {
		t.Error("illegal length")
	}

}

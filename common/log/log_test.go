package log

import "testing"

func TestDebugJson(t *testing.T) {
	type testedStruct struct {
		Foo string `json:"foo"`
		Bar string `json:"-"`
	}
	ts := &testedStruct{
		Foo: "foo",
		Bar: "bar",
	}
	var err error
	if err = DebugJson("test", ts); err != nil {
		t.Error("debug json err:", err)
	}
	if err = DebugJson("test", *ts); err != nil {
		t.Error("debug json err:", err)
	}
	if err = DebugJson(123, ts); err == nil {
		t.Error("debug json should occur err")
	}
}

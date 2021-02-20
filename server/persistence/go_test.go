package persistence

import (
	"testing"
)

type b struct {
	i int
}
func a() interface{} {
	return nil
}
func TestRDBInit(t *testing.T) {
	if e := a(); e==nil {
		t.Log("h")
	}
	//t.Logf("%+v\n",a().(b))
}

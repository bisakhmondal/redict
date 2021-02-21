package persistence

import (
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"
)

func TestRDB_LoadDump(t *testing.T) {
	tt := Transaction{
		Uid:   1,
		Key:   "bisakh",
		Value: "Mondal",
	}
	t1 := Transaction{
		Uid:   2,
		Key:   "Mondal",
		Value: "hello",
	}

	f, _ := os.Create("a.json")
	r, _ := Marshal(tt)
	io.Copy(f,r)
	rr,_ := Marshal(t1)
	io.Copy(f, rr)
	f.Close()

	ff ,_ := os.Open("a.json")
	dec := json.NewDecoder(ff)

	defer ff.Close()
	var ts Transaction
	for {
		err := dec.Decode(&ts)
		if err!=nil{
			return
		}
		t.Logf("%+v", ts)

	}
	//dec.Decode(&ts)
	//t.Logf("%+v", ts)

}

func TestNewRDB(t *testing.T) {
	rdb := NewRDB()

	rdb.LoadDump("../2021-02-21T12:40:35+05:30.json")
	t.Log(len(rdb.queue.Items))
	time.Sleep(5*time.Second)
	for {
		item := rdb.queue.Pop()
		if item ==nil{
			return
		}

		t.Logf("%+v\n", item.(Transaction))
	}
}

package persistence

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

type RDB struct {
	filename string
	ticker *time.Ticker //snapshotting period
	queue *Queue //thread safe queue
	quit chan struct{}
}

func NewRDB() *RDB{
	return &RDB{
		filename: "snapshots/"+time.Now().Format(time.RFC3339)+".json",
		ticker:   time.NewTicker(time.Second),
		queue:    newQueue(),
		quit: make(chan struct{}),
	}
}

//Get Queue
func (r * RDB)GetQueue() *Queue{
	return r.queue
}

//Initialize for fresh Dump
func RDBInit () * RDB{
	rdb := NewRDB()
	var err error
	f , err := os.Create(rdb.filename)
	defer f.Close()

	if err!= nil {
		panic(err)
	}

	go func() {
		log.Println("Concurrent Backup Initiated")
		for{
			select {
			case <- rdb.ticker.C:
				rdb.periodicDumpContent()
			case <- rdb.quit:
				log.Println("Backup successful")
				return

			}
		}
	}()
	
	return rdb
}

//API for appending each individual Transaction
func (rdb *RDB) WriteTransaction(uid int, key, value string){

	rdb.queue.Push(Transaction{
		Uid: uid,
		Key: key,
		Value: value,
	})
}

//function for write transactions in JSON during snapshotting
func (rdb * RDB)periodicDumpContent(){
	fid, err := os.OpenFile(rdb.filename , os.O_WRONLY|os.O_APPEND, 0600)

	if err!= nil{
		panic(err)
	}

	for {
		out := rdb.queue.Pop()
		if out == nil {
			return
		}

		data, err := Marshal(out.(Transaction))

		if err != nil {
			log.Printf("Marshalling error %+v\n", out.(Transaction))
			return
		}
		io.Copy(fid, data)
	}
}

//Load previous Storage Dump
func (rdb *RDB) LoadDump(filename string) {
	var err error
	fid , err := os.Open(filename)
	defer fid.Close()

	if err != nil{
		panic(err)
	}

	var t Transaction
	dec := json.NewDecoder(fid)

	for {
		err = dec.Decode(&t)
		if err != nil {
			return
		}
		rdb.WriteTransaction(t.Uid, t.Key, t.Value)
	}
}

//Function to shut down incremental Snapshotting
func (rdb * RDB) Quit(){
	rdb.ticker.Stop()
	close(rdb.quit)
}

package persistence

import (
	"compress/gzip"
	"encoding/gob"
	"log"
	"os"
	"time"
)

type RDB struct {
	filename string
	ticker *time.Ticker //snapshotting period
	queue *Queue //thread safe queue
	quit chan struct{}
	fid *os.File
}

func NewRDB() *RDB{
	return &RDB{
		filename: time.Now().Format(time.RFC3339)+".log",
		ticker:   time.NewTicker(time.Second),
		queue:    newQueue(),
		quit: make(chan struct{}),
	}
}

//Initialize for fresh Dump
func RDBInit () * RDB{
	rdb := NewRDB()
	var err error
	rdb.fid, err = os.Create(rdb.filename)

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

//API for appending each individual transaction
func (rdb *RDB) WriteTransaction(uid int, key, value string){

	rdb.queue.Push(transaction{
		Uid: uid,
		Key: key,
		Value: value,
	})
}

//function for write GZIP compressed transactions in file during snapshotting
func (rdb * RDB)periodicDumpContent(){
	fz := gzip.NewWriter(rdb.fid)
	defer fz.Close()

	enc := gob.NewEncoder(fz)

	for {
		out := rdb.queue.Pop()
		if out == nil {
			return
		}

		err := enc.Encode(out)
		if err != nil {
			log.Printf("encoding error %+v\n", out.(transaction))
			return
		}
	}
}

//Load previous Storage Dump
func (rdb *RDB) LoadDump(filename string) {
	var err error
	rdb.fid, err = os.Open(filename)
	defer rdb.fid.Close()

	if err != nil{
		panic(err)
	}

	fz, err := gzip.NewReader(rdb.fid)
	defer fz.Close()
	if err!= nil {
		panic(err)
	}

	dec := gob.NewDecoder(fz)

	var t transaction

	for res := dec.Decode(&t); res!= nil; {
		rdb.WriteTransaction(t.Uid, t.Key, t.Value)
	}
}

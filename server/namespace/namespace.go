package namespace

import "redict/server/persistence"

//blank symbol
const blank = "<blank>"

type Namespace struct {
	uid int
	dict map[string]string
	rdb * persistence.RDB
}

func New(uid int)*Namespace{
	return &Namespace{
		uid: uid,
		dict: map[string]string{},
	}
}

func (n * Namespace) Put(key, value string){
	n.dict[key]  = value

	//for in memory persistence
	n.rdb.WriteTransaction(n.uid, key, value)
}

func (n* Namespace) Get(key string) string{
	val, ok := n.dict[key]
	if !ok {
		return blank
	}

	return val
}


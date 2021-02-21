package namespace

import (
	"fmt"
	"log"
	"redict/server/persistence"
	"sync"
)

type Container struct {
	sync.Mutex
	list map[int]*Namespace
}

func NewContainer() *Container{
	return &Container{
		list: map[int]*Namespace{},
	}
}

func (c * Container) Push(n* Namespace){
	c.Lock()
	defer c .Unlock()
	c.list[n.uid] = n
}

func (c* Container) Delete (n *Namespace){
	c.Lock()
	defer c.Unlock()

	delete(c.list, n.uid)

}

//if Manager wants to demote its role to a particular User
func (c* Container) GetNamespace(uid int) *Namespace{
	namesp, ok := c.list[uid]
	if !ok	{
		return nil
	}
	return namesp
}

//return current stats to the Manager
func (c* Container) GetStats() string{
	info := ""

	//to avoid dict change during lookup
	c.Lock()
	for key, value  := range c.list{
		info += fmt.Sprintf("User: %d | Total kv store %d\n", key, len(value.dict))
	}
	c.Unlock()

	if info=="" {
		return "Not a single client ever attached\n"
	}
	return info
}

//Meant for Manager
//return all values of Clients KV store for a particular key
func (c * Container)Get(key string) string{
	info := "key: " + key + "\n"
	info += "Found at-------------------------\n"

	c.Lock()
	for _, namespace := range c.list {
		val := namespace.Get(key)

		if val != blank {
			info += fmt.Sprintf("User: %d | value: %s\n", namespace.uid, val)
		}
	}
	c.Unlock()

	return info
}

//Meant for Manager
//Override kv pair of all clients Namespace
func (c* Container)Put(key, value string){

	c.Lock()
	for _, namespace := range c.list {
		namespace.Put(key, value)
	}
	c.Unlock()
}

//Initalize
//Function to restore previous backup
func (c * Container)Restore(q *persistence.Queue, rdb * persistence.RDB) int{
	lastClient := -1

	for {
		entry := q.Pop()

		if entry==nil{
			break
		}

		item := entry.(persistence.Transaction)
		_ , ok := c.list[item.Uid]
		if !ok {
			c.Push(New(item.Uid, rdb))
		}
		c.list[item.Uid].Put(item.Key, item.Value)

		//update total client attached
		if item.Uid > lastClient {
			lastClient = item.Uid
		}
	}
	log.Println("Previously Attched Clients", lastClient+1)

	return lastClient + 1
}
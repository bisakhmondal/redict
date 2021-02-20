package persistence

import "sync"

const capacity = 256

type Queue struct {
	sync.Mutex
	Items []interface{}
}

func newQueue() *Queue{
	return &Queue{
		Items: make([]interface{}, 0, capacity),
	}
}

func (q *Queue) Push(item interface{}) {
	q.Lock()
	defer q.Unlock()
	q.Items = append(q.Items, item)
}

func (q *Queue) Pop() interface{} {
	q.Lock()
	defer q.Unlock()
	if len(q.Items) == 0 {
		return nil
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item
}

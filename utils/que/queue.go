package que

import (
	"sync"
)

type Queue struct {
	lock   *sync.Mutex
	Values []interface{}
}

func Init() *Queue {
	return &Queue{&sync.Mutex{}, make([]interface{}, 0)}
}

func (q *Queue) Enqueue(x interface{}) {
	for {
		q.lock.Lock()
		q.Values = append(q.Values, x)
		q.lock.Unlock()
		return
	}
}

func (q *Queue) Dequeue() interface{} {
	for {
		if len(q.Values) > 0 {
			q.lock.Lock()
			x := q.Values[0]
			q.Values = q.Values[1:]
			q.lock.Unlock()
			return x
		}
		return nil
	}
	return nil
}

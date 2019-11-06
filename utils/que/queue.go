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
	q.lock.Lock()
	defer q.lock.Unlock()
	for {
		q.Values = append(q.Values, x)
		return
	}
}

func (q *Queue) Dequeue() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	for {
		if len(q.Values) > 0 {
			x := q.Values[0]
			q.Values = q.Values[1:]
			return x
		}
		return nil
	}
	return nil
}

func (q *Queue) Get(index int) interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	for {
		if len(q.Values) > index {
			return q.Values[index]
		}
		return nil
	}
	return nil
}

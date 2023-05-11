package gcache

import "sync"

type lruLink[V int | int64 | float64 | string] struct {
	head *node[V]
	sync.RWMutex
}

type node[V int | int64 | float64 | string] struct {
	v    V
	next *node[V]
}

func (l *lruLink[V]) flush(v V) {
	l.Lock()
	defer l.Unlock()

	n := &node[V]{
		v:    v,
		next: nil,
	}

	if l.head == nil {
		l.head = n
	} else {
		if l.head.v == v {
			if l.head.next == nil {
				return
			} else {
				l.head = l.head.next
			}
		}

		cur := l.head
		for cur.next != nil {
			if cur.next.v == v {
				cur.next = cur.next.next
			} else {
				cur = cur.next
			}
		}
		cur.next = n
	}
}

func (l *lruLink[V]) popHead() (V, bool) {
	l.RLock()
	defer l.RUnlock()
	var v V
	if l.head == nil {
		return v, false
	}
	v = l.head.v
	l.head = l.head.next
	return v, true
}

func (l *lruLink[V]) remove(v V) {
	l.Lock()
	defer l.Unlock()
	if l.head.v == v {
		l.head = l.head.next
		return
	}
	curr := l.head
	for curr.next != nil {
		if curr.next.v == v {
			curr.next = curr.next.next
			return
		} else {
			curr = curr.next
		}
	}
}

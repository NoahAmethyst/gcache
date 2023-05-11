package gcache

import "sync"

// lruLink a linklist that save active keys sorted by activity degree.
// The more in front, the longer that  haven't been used.
type lruLink[V int | int64 | float64 | string] struct {
	head *node[V]
	sync.RWMutex
}

type node[V int | int64 | float64 | string] struct {
	v    V
	next *node[V]
}

// flush add new value in the list tail or switch old key from the last index to tail
// Make the list sorted by active degree and all values are unique in the list
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

// popHead return the head value (longest key haven't been used) and delete it
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

// remove remove specific value from list
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

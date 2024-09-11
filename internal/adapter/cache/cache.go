package cache

import (
	"sync"
)

type Lru[K comparable, V any] struct {
	capacity int
	head     *Element[K, V]
	tail     *Element[K, V]
	data     map[K]*Element[K, V]
	mu       sync.RWMutex
}

type Element[K comparable, V any] struct {
	k    K
	v    V
	prev *Element[K, V]
	next *Element[K, V]
}

func NewLru[K comparable, V any](capacity int) *Lru[K, V] {
	return &Lru[K, V]{
		capacity: capacity,
		data:     make(map[K]*Element[K, V]),
	}
}

func (l *Lru[K, V]) Get(key K) (v *V, exists bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var element *Element[K, V]

	element, exists = l.data[key]

	if !exists {
		return
	}

	l.remove(element)
	l.moveToHead(element)

	return &element.v, true
}

func (l *Lru[K, V]) Put(key K, value V) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if element, exists := l.data[key]; exists {
		element.v = value
		l.remove(element)
		l.moveToHead(element)
	} else {
		element = &Element[K, V]{
			k: key,
			v: value,
		}

		if len(l.data) == 0 {
			l.head = element
			l.tail = element

			l.data[key] = element

			return
		}

		if len(l.data) >= l.capacity {
			delete(l.data, l.tail.k)
			l.remove(l.tail)
		}

		l.moveToHead(element)
		l.data[key] = element
	}
}

func (l *Lru[K, V]) moveToHead(element *Element[K, V]) {
	element.prev = nil
	element.next = l.head
	l.head.prev = element
	l.head = element

}

func (l *Lru[K, V]) remove(element *Element[K, V]) {
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}

	if element == l.head {
		if l.head.next != nil {
			l.head = l.head.next
			l.head.prev = nil
		}
	}

	if element == l.tail {
		if l.tail.prev != nil {
			l.tail = l.tail.prev
			l.tail.next = nil
		}
	}
}

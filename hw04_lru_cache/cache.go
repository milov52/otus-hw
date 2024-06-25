package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	e, ok := l.items[key]
	if ok {
		l.items[key].Value = value
		l.queue.MoveToFront(e)
		return true
	}

	if l.queue.Len() == l.capacity {
		for k, v := range l.items {
			if v == l.queue.Back() {
				delete(l.items, k)
				l.queue.Remove(v)
			}
		}
	}
	l.items[key] = l.queue.PushFront(value)
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	e, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(e)
		return e.Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	NewCache(l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

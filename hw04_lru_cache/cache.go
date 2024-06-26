package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheEntry struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	item, ok := l.items[key]
	if ok {
		item.Value = cacheEntry{key, value}
		l.queue.MoveToFront(item)
		return true
	}

	if l.queue.Len() == l.capacity {
		back := l.queue.Back()
		if back != nil {
			entry := back.Value.(cacheEntry)
			delete(l.items, entry.key)
			l.queue.Remove(back)
		}
	}

	l.items[key] = l.queue.PushFront(cacheEntry{key, value})
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := l.items[key]
	if ok {
		l.queue.MoveToFront(item)
		return item.Value.(cacheEntry).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.front == nil {
		l.front = &ListItem{Value: v}
		l.back = l.front
	} else {
		l.front = &ListItem{Value: v, Next: l.front}
		l.front.Next.Prev = l.front
	}
	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.front == nil {
		l.front = &ListItem{Value: v}
		l.back = l.front
	} else {
		l.back = &ListItem{Value: v, Prev: l.back}
		l.back.Prev.Next = l.back
	}
	l.len++
	return l.front
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}

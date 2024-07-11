package main

import (
	"fmt"
	hw04_lru_cache "github.com/milov52/otus-hw/hw04_lru_cache"
)

func main() {
	l := hw04_lru_cache.NewList()

	l.PushFront(10) // [10]
	l.PushBack(20)  // [10, 20]
	l.PushBack(30)  // [10, 20, 30]

	middle := l.Front().Next // 20
	l.Remove(middle)         // [10, 30]

	for i, v := range [...]int{40, 50, 60, 70, 80} {
		if i%2 == 0 {
			l.PushFront(v)
		} else {
			l.PushBack(v)
		}
	}
	fmt.Println(l)
}

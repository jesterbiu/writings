package main

import "fmt"

type Iterator[T any] interface {
	Next() bool
	Get() T
}

type Pair[T1, T2 any] struct {
	first  T1
	second T2
}

type sliceIter[T any] struct {
	s []T
	i int // start at off-by-1
}

func (iter *sliceIter[T]) Next() bool {
	iter.i++
	return iter.i < len(iter.s)
}

func (iter *sliceIter[T]) Get() T {
	return iter.s[iter.i]
}

func SliceIter[T any]() Iterator[T] {
	return nil
}

type mapIter[T any] struct{}

func main() {
	fmt.Println("Hello, 世界")
}

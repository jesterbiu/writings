package main

import (
	"fmt"
	"sync"
	"testing"
)

func newThread(wg *sync.WaitGroup, f func()) {
	go func() {
		defer wg.Done()
		f()
	}()
}

func waitGroupN(n int) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(n)
	return wg
}

// go test -run Test_NaiveSingleton -v -race
func Test_NaiveSingleton(t *testing.T) {
	type singleton struct {
		name string
	}
	var instance *singleton
	makeSingleton := func() *singleton {
		fmt.Println("makeSingleton")
		return new(singleton)
	}
	getSingleton := func() *singleton {
		if instance == nil {
			instance = makeSingleton()
		}
		return instance
	}
	n := 10
	wg := waitGroupN(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i, getSingleton())
		}(i)
	}
	wg.Wait()
}

// go test -run Test_MutexSingleton -v -race
func Test_MutexSingleton(t *testing.T) {
	type singleton struct {
		name string
	}
	var (
		instance *singleton
		mutex    sync.Mutex
	)
	makeSingleton := func() *singleton {
		fmt.Println("makeSingleton")
		return new(singleton)
	}
	getSingleton := func() *singleton {
		mutex.Lock()
		defer mutex.Unlock()
		if instance == nil {
			instance = makeSingleton()
		}
		return instance
	}
	n := 10
	wg := waitGroupN(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i, getSingleton())
		}(i)
	}
	wg.Wait()
}

// go test -run Test_NaiveDBL -v -race
func Test_NaiveDBL(t *testing.T) {
	type singleton struct {
		name string
	}
	var (
		instance *singleton
		mutex    sync.Mutex
	)
	makeSingleton := func() *singleton {
		fmt.Println("makeSingleton")
		return new(singleton)
	}
	getSingleton := func() *singleton {
		if instance == nil {
			mutex.Lock()
			defer mutex.Unlock()
			if instance == nil {
				instance = makeSingleton()
			}
		}
		return instance
	}
	n := 100
	wg := waitGroupN(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i, getSingleton())
		}(i)
	}
	wg.Wait()
}

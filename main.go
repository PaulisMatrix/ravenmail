package main

import (
	"errors"
	"fmt"
	"os"
)

type StoreIface interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

var _ StoreIface = (*Store)(nil)

func GetNewKV(capacity int) *Store {
	s := &Store{
		KVMap:    make(map[string]interface{}, capacity),
		FIFO:     getDLL(),
		capacity: capacity,
	}
	return s
}

func (s *Store) Set(key string, value interface{}) error {
	// check the cur len > capacity, delete the head node.
	fmt.Println("current capacity", s.FIFO.capacity)

	// add node to the DLL
	if s.FIFO.capacity >= s.capacity {
		// evict the head node and update the capacity
		fmt.Println("capacity breached, deleting head node...")
		// how to get the key
		node := s.FIFO.deleteNode()
		delete(s.KVMap, node.key)
	}

	newNode := s.FIFO.addNode(key, value)
	// key -> nodeRef
	s.KVMap[key] = newNode

	return nil
}

func (s *Store) Get(key string) (interface{}, error) {
	nodeRef, ok := s.KVMap[key]

	if !ok {
		return nil, errors.New("Key doesn't exist")
	}

	node := nodeRef.(*Node)

	return node.val, nil
}

type Store struct {
	// store the key value pairs
	KVMap    map[string]interface{}
	FIFO     *DoublyLinkedList
	capacity int
	StoreIface
}

func main() {
	s := GetNewKV(1)

	s.Set("hello", "world")
	s.Set("first", 100)

	val, err := s.Get("hello")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("value received: ", val)
}
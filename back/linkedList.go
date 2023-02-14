package main

import (
	"log"
)

type LinkedList struct {
	head *Node
}

type Node struct {
	data Monitor
	next *Node
}

func (l *LinkedList) add(data Monitor) {
	node := l.head

	if node == nil {
		l.head = &Node{data: data}
		return
	}

	for node.next != nil {
		node = node.next
	}

	newNode := &Node{data: data}
	node.next = newNode
}

func (l *LinkedList) len() int {
	node := l.head
	length := 0
	for node != nil {
		node = node.next
		length += 1
	}
	return length
}

func (l *LinkedList) toArray() []*Monitor {
	node := l.head
	array := make([]*Monitor, l.len())
	index := 0
	for node != nil {
		array[index] = &node.data
		node = node.next
		index++
	}
	return array
}

func (l *LinkedList) get(i int) *Monitor {
	node := l.head
	if i >= l.len() {
		log.Println("Index out of bounds")
		return nil
	}
	for j := 0; j < i; j++ {
		node = node.next
	}
	return &node.data
}

func (l *LinkedList) set(i int, data Monitor) {
	node := l.head
	if i >= l.len() {
		log.Println("Index out of bounds")
		return
	}
	for j := 0; j < i; j++ {
		node = node.next
	}
	node.data = data
}

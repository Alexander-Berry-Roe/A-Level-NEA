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

func (l *LinkedList) remove(data Monitor) {
	node := l.head
	var prev *Node

	// Traverse the list to find the node to remove
	for node != nil {
		if node.data == data {
			// Found the node to remove

			if prev == nil {
				// Removing the head of the list
				l.head = node.next
			} else {
				// Removing a node in the middle or end of the list
				prev.next = node.next
			}
			return
		}

		// Keep track of the previous node
		prev = node
		node = node.next
	}
}

package queue

import (
	"fmt"
	"sync"
)

type node struct {
	data interface{}
	next *node
}

// Thread safe FIFO Queue data structure
type Queue struct {
	head, tail  *node
	size, items int64
	mutex       *sync.RWMutex
}

// Create a new queue
// If the parameter size is passed, the queue will be of that fixed size
func New(size ...int64) *Queue {
	queue := &Queue{mutex: new(sync.RWMutex)}
	queue.size = -1
	if len(size) > 0 {
		queue.size = size[0]
	}
	return queue
}

// Return the Queue max capacity (size)
// NOTE: -1 is returned for unlimited Queues
func (queue *Queue) Cap() int64 {
	queue.mutex.RLock()
	defer queue.mutex.RUnlock()

	return queue.size
}

// Return the Queue current length
func (queue *Queue) Len() int64 {
	queue.mutex.RLock()
	defer queue.mutex.RUnlock()

	return queue.items
}

// Push an element at the end of the Queue
func (queue *Queue) Push(element interface{}) error {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if queue.items == queue.size {
		return fmt.Errorf(
			"Can't push %v, Queue beyond limits (%d)", element, queue.size)
	}

	n := &node{data: element}
	if queue.tail == nil {
		queue.firstAndLatest(n)
	} else {
		queue.pushToEnd(n)
	}
	queue.increment()
	return nil
}

// 从队列末尾提取元素
func (queue *Queue) Pop() interface{} {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if queue.items == 0 {
		return nil
	}

	// move the head forward
	data := queue.head.data
	n := queue.head
	queue.head = n.next
	n = nil // free the memory
	queue.items--
	queue.checkEmptyness()

	return data
}

// Get all the values from the Queue as an interface{} slice
func (queue *Queue) Values() []interface{} {
	queue.mutex.RLock()
	defer queue.mutex.RUnlock()

	data := make([]interface{}, queue.items)
	n := queue.head
	if n == nil {
		return data
	}
	data[0] = n.data

	i := 1
	for {
		next := n.next
		if next == nil {
			break
		}
		data[i] = next.data
		n = next
		i++
	}

	return data
}

// push an element as first and last element of the queue
func (queue *Queue) firstAndLatest(n *node) {
	queue.tail = n
	queue.head = n
}

// push an element to the end of the queue and move the prev tail.next pointer
func (queue *Queue) pushToEnd(n *node) {
	queue.tail.next = n
	queue.tail = n
}

// increment the items counter on the queue
func (queue *Queue) increment() {
	queue.items++
}

// check if the queue is empty, if so, set the head and the tail to nil
func (queue *Queue) checkEmptyness() {
	if queue.items == 0 {
		queue.head = nil
		queue.tail = nil
	}
}

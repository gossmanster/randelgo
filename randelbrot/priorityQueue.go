package randelbrot

import (
	"container/list"
)

type item struct {
	set        *MandelbrotSet
	evaluation float64
}

type priorityQueue struct {
	maxSize int
	items   list.List
}


func newPriorityQueue(maxSize int) *priorityQueue {
	q := new(priorityQueue)
	q.maxSize = maxSize
	q.items.Init()
	
	return q
}

func (q *priorityQueue) len() int {
	return q.items.Len()
}

func (q *priorityQueue) push(set *MandelbrotSet, evaluation float64) {
	newItem := new(item)
	newItem.set = set
	newItem.evaluation = evaluation
	
	// Add in order of priority, keeping the list sorted
	added := false
	for e := q.items.Front(); e != nil; e = e.Next() {
		oldItem := e.Value.(*item)
		if (newItem.evaluation > oldItem.evaluation) {
			q.items.InsertBefore(newItem, e)
			added = true
			break
		}		
	}	
	if (!added){
		q.items.PushBack(newItem)
	}
	
	// Once the list gets filled up, remove the last item in priority
	if (q.len() > q.maxSize) {
		q.items.Remove(q.items.Back())
	}
}

func (q *priorityQueue) pop() *MandelbrotSet {
	top := q.items.Front()
	q.items.Remove(top)
	
	return top.Value.(*item).set
}

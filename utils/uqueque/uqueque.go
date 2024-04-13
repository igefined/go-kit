package uqueque

type item struct {
	value any
	next  *item
}

type Queue struct {
	head *item
	tail *item
}

func NewQueue() *Queue {
	dummy := &item{}
	return &Queue{
		head: dummy,
		tail: dummy,
	}
}

func (q *Queue) Push(value any) {
	q.tail.next = &item{value: value}
	q.tail = q.tail.next
}

func (q *Queue) Pop() any { //nolint:ireturn
	if q.head.value == q.tail.value {
		return nil
	}

	value := q.head.next.value
	q.head = q.head.next
	return value
}

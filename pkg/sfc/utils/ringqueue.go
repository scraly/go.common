package utils

// Ringqueue struct
type Ringqueue struct {
	nodes []interface{}
	head  int
	tail  int
	cnt   int
}

// NewRingqueue constructs a Ringqueue
func NewRingqueue() *Ringqueue {
	return &Ringqueue{
		nodes: make([]interface{}, 2),
	}
}

func (q *Ringqueue) resize(n int) {
	nodes := make([]interface{}, n)
	if q.head < q.tail {
		copy(nodes, q.nodes[q.head:q.tail])
	} else {
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.tail])
	}

	q.tail = q.cnt % n
	q.head = 0
	q.nodes = nodes
}

// Add object to the queue (last position)
func (q *Ringqueue) Add(i interface{}) {
	if q.cnt == len(q.nodes) {
		// Also tested a grow rate of 1.5, see: http://stackoverflow.com/questions/2269063/buffer-growth-strategy
		// In Go this resulted in a higher memory usage.
		q.resize(q.cnt * 2)
	}
	q.nodes[q.tail] = i
	q.tail = (q.tail + 1) % len(q.nodes)
	q.cnt++
}

// Remove head object of the queue
func (q *Ringqueue) Remove() (interface{}, bool) {
	if q.cnt == 0 {
		return 0, false
	}
	i := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.cnt--

	if n := len(q.nodes) / 2; n > 2 && q.cnt <= n {
		q.resize(n)
	}

	return i, true
}

// Cap function
func (q *Ringqueue) Cap() int {
	return cap(q.nodes)
}

// Len function
func (q *Ringqueue) Len() int {
	return q.cnt
}

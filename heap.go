package orderedconcurrently

type processInputHeap []*processInput

func (h processInputHeap) Len() int {
	return len(h)
}

func (h processInputHeap) Less(i, j int) bool {
	return h[i].order < h[j].order
}

func (h processInputHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *processInputHeap) Push(x interface{}) {
	*h = append(*h, x.(*processInput))
}

func (h processInputHeap) Peek() interface{} {
	if h.Len() > 0 {
		return h[0]
	}
	return nil
}

func (h *processInputHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

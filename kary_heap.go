package cron

func NewKaryHeap(k int, cmp func(i, j *CronJob) bool) *KaryHeap {
	h := &KaryHeap{
		k: k,

		cmp: cmp,
	}

	return h
}

type KaryHeap struct {
	array []*CronJob
	k     int
	cmp   func(i, j *CronJob) bool // i is root if true

}

func (h *KaryHeap) swap(i, j int) {
	h.array[i], h.array[j] = h.array[j], h.array[i]
}

func (h *KaryHeap) RestoreDown(idx int) {
	if idx < 0 {
		return
	}

	if idx+1 >= len(h.array) {
		return
	}
	h.down(idx)
}

func (h *KaryHeap) down(idx int) {
	var (
		ln = len(h.array)
		k  = h.k
	)

	ci, pc := idx, idx
	for ci < ln {
		cl := ci*k + k
		for i := cl - k + 1; i <= cl && i < ln; i++ {
			if h.cmp(h.array[i], h.array[pc]) {
				pc = i
			}
		}
		if ci != pc {
			h.swap(ci, pc)
			ci = pc
		} else {
			break
		}
	}
}

func (h *KaryHeap) RestoreUp(idx int) {
	if idx <= 0 {
		return
	}
	if idx > len(h.array)-1 {
		return
	}
	h.up(idx)
}

func (h *KaryHeap) up(idx int) {
	k := h.k
	ci := idx
	parent := (ci - 1) / k
	for ci > 0 {
		if h.cmp(h.array[ci], h.array[parent]) {
			h.swap(ci, parent)
			ci = parent
			parent = (ci - 1) / k
		} else {
			break
		}
	}
}

func (h *KaryHeap) Push(elem *CronJob) {
	h.push(elem)
}

func (h *KaryHeap) push(elem *CronJob) {
	h.array = append(h.array, elem)
	h.up(len(h.array) - 1)
}

func (h *KaryHeap) Pop() *CronJob {
	return h.pop()
}

func (h *KaryHeap) pop() *CronJob {
	if len(h.array) == 0 {
		return nil
	}
	elem := h.array[0]
	h.swap(0, len(h.array)-1)
	h.array = h.array[:len(h.array)-1]
	h.down(0)
	return elem
}

func (h *KaryHeap) Peek(idx int) *CronJob {
	return h.peek(idx)
}

func (h *KaryHeap) peek(idx int) *CronJob {
	if idx < 0 || idx >= len(h.array) {
		return nil
	}
	return h.array[idx]
}

func (h *KaryHeap) Walk(walk func(*CronJob) bool) int {
	for k, v := range h.array {
		if walk(v) {
			return k
		}
	}
	return -1
}

func (h *KaryHeap) Remove(idx int) *CronJob {
	return h.remove(idx)
}

func (h *KaryHeap) remove(idx int) *CronJob {
	ret := h.array[idx]
	ln := len(h.array)
	h.array[idx] = h.array[ln-1]
	h.array = h.array[:ln-1]

	if h.cmp(h.array[idx], h.array[(idx-1)/h.k]) {
		h.up(idx)
	} else {
		h.down(idx)
	}

	return ret
}

func (h *KaryHeap) Len() int {
	return len(h.array)
}

func (h *KaryHeap) BuildHeap() {
	ln := len(h.array)
	for i := (ln - 1) / h.k; i >= 0; i-- {
		h.down(i)
	}
}

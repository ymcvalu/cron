package cron

import "testing"

type T struct {
	i    int
	s, d []interface{}
}

func TestDown(t *testing.T) {
	var set = []T{
		T{
			1,
			[]interface{}{1, 8, 5, 7, 3, 9, 10, 11, 12, 13, 14, 15},
			[]interface{}{1, 3, 5, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		},
	}

	h := NewKaryHeap(func(i, j interface{}) bool {
		ii, ij := i.(int), j.(int)
		return ii < ij
	}, Kary(3))
	for _, e := range set {
		h.array = e.s
		h.RestoreDown(e.i)
		for i := range h.array {
			if h.array[i] != e.d[i] {
				t.Error(e.s, e.d)
				break
			}
		}
	}
}

func TestUp(t *testing.T) {
	var set = []T{
		T{
			3,
			[]interface{}{3, 8, 5, 7},
			[]interface{}{3, 7, 5, 8},
		},
	}

	h := NewKaryHeap(func(i, j interface{}) bool {
		ii, ij := i.(int), j.(int)
		return ii < ij
	})
	for _, e := range set {
		h.array = e.s
		h.RestoreUp(e.i)
		for i := range h.array {
			if h.array[i] != e.d[i] {
				t.Error(e.s, e.d)
				break
			}
		}
	}
}

func TestPush(t *testing.T) {
	var set = []T{
		T{
			3,
			[]interface{}{8, 7},
			[]interface{}{7, 8},
		},
		T{
			3,
			[]interface{}{8, 7, 5, 3},
			[]interface{}{3, 5, 7, 8},
		},
	}
	h := NewKaryHeap(func(i, j interface{}) bool {
		ii, ij := i.(int), j.(int)
		return ii < ij
	})

	for _, e := range set {
		h.array = []interface{}{}
		for _, v := range e.s {
			h.Push(v)
		}
		for i := range h.array {
			if h.array[i] != e.d[i] {
				t.Error(h.array, e.d)
				break
			}
		}
	}
}

func TestPop(t *testing.T) {
	var set = []T{

		T{
			3,
			[]interface{}{3, 9, 7, 10, 12, 8, 9},
			[]interface{}{3, 7, 8, 9, 9, 10, 12},
		},
	}
	h := NewKaryHeap(func(i, j interface{}) bool {
		ii, ij := i.(int), j.(int)
		return ii < ij
	})

	for _, e := range set {
		h.array = e.s
		arr := []interface{}{}
		for {
			elem := h.Pop()
			if elem != nil {
				arr = append(arr, elem)
			} else {
				break
			}
		}
		for i := range arr {
			if arr[i] != e.d[i] {
				t.Error(arr, e.d)
				break
			}
		}
	}
}

func TestRemove(t *testing.T) {
	var set = []T{

		T{
			3,
			[]interface{}{1, 3, 5, 7, 9, 6, 8},
			[]interface{}{1, 7, 5, 8, 9, 6},
		},
	}
	h := NewKaryHeap(func(i, j interface{}) bool {
		ii, ij := i.(int), j.(int)
		return ii < ij
	})

	for _, e := range set {
		h.array = e.s

		h.Remove(1)
		for i := range h.array {
			if h.array[i] != e.d[i] {
				t.Error(h.array, e.d)
				break
			}
		}
	}
}

func TestPeek(t *testing.T) {
	var set = []T{

		T{
			3,
			[]interface{}{1, 3, 5, 7, 9, 6, 8},
			[]interface{}{1, 7, 5, 8, 9, 6},
		},
	}
	h := NewKaryHeap(func(i, j interface{}) bool {
		ii, ij := i.(int), j.(int)
		return ii < ij
	})

	for _, e := range set {
		h.array = e.s
		h.Remove(1)
		for i := range h.array {
			if h.Peek(i) != e.d[i] {
				t.Error(h.array, e.d)
				break
			}
		}

	}
}

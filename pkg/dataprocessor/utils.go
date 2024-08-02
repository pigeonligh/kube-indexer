package dataprocessor

type sliceRange struct {
	start, end *int
}

func (r sliceRange) maintain(length int) (int, int) {
	var f, t int
	if r.start == nil {
		f = 0
	} else {
		f = *r.start
	}
	if r.end == nil {
		t = length
	} else {
		t = *r.end
	}

	if f < 0 {
		f = 0
	}
	if t > length {
		t = length
	}
	return f, t
}

func rangeT[T any](arr []T, ranges ...sliceRange) []T {
	ret := make([]T, 0)
	for _, r := range ranges {
		f, t := r.maintain(len(arr))
		for i := f; i < t; i++ {
			ret = append(ret, arr[i])
		}
	}
	return ret
}

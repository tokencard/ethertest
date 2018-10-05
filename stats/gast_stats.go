package stats

import "sort"

// Uint64Min calculates minimum value in the array.
// Returns 0 if array is empty.
func Uint64Min(data []uint64) uint64 {
	if len(data) == 0 {
		return 0
	}
	min := data[0]
	for _, v := range data[1:] {
		if min > v {
			min = v
		}
	}
	return min
}

// Uint64Max calculates maximum value in the array.
// Returns 0 if array is empty.
func Uint64Max(data []uint64) uint64 {
	if len(data) == 0 {
		return 0
	}
	max := data[0]
	for _, v := range data[1:] {
		if max < v {
			max = v
		}
	}
	return max
}

// Uint64Median calculates median value in the array.
// Returns 0 if array is empty.
func Uint64Median(data []uint64) uint64 {
	if len(data) == 0 {
		return 0
	}

	c := make([]uint64, len(data))
	copy(c, data)

	sort.Slice(c, func(i int, j int) bool {
		return c[i] < c[j]
	})

	return c[len(c)/2]
}

type sortableUint64Slice []uint64

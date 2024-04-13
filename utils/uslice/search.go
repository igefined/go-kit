package uslice

type Number interface {
	int | int64 | uint64 | float64 | float32
}

func Search[T Number](nums []T, elem T) int {
	low := 0
	high := len(nums) - 1

	for low <= high {
		mid := (low + high) / 2
		guess := nums[mid]
		if elem == guess {
			return mid
		}

		if elem > guess {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

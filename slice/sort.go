package slice

func Sort[T Number](nums []T) []T {
	if len(nums) < 2 {
		return nums
	}

	low := 0
	high := len(nums) - 1

	qsort[T](nums, low, high)
	
	return nums
}

func partition[T Number](nums []T, low, high int) int {
	pivot := nums[high]

	j := low - 1

	for i := low; i < high; i++ {
		if nums[i] < pivot {
			j++
			nums[i], nums[j] = nums[j], nums[i]
		}
	}
	
	nums[j + 1], nums[high] = pivot, nums[j + 1]
	return j + 1
}

func qsort[T Number](nums []T, low, high int) {
	if low < high{
		pivot := partition[T](nums, low, high)

		qsort[T](nums, low, pivot - 1)
		qsort[T](nums, pivot + 1, high)
	}
}
package sort

type Quick struct{}

func (q Quick) Sort(arr []int) []int {
	shuffle(arr)
	q.sort(arr, 0, len(arr)-1)
	return arr
}

func (q Quick) sort(arr []int, lo, hi int) {
	if lo >= hi {
		return
	}

	p := q.partition(arr, lo, hi)

	q.sort(arr, lo, p-1)
	q.sort(arr, p+1, hi)
}

func (q Quick) partition(arr []int, lo, hi int) int {
	pivot := arr[lo]
	i, j := lo+1, hi

	for i <= j {
		// 以pivot为中间值，找到arr的左边大于等于pivot，且右边小于pivot
		// 并进行交换
		for i < hi && arr[i] <= pivot {
			i++
		}

		for j > lo && arr[j] > pivot {
			j--
		}

		if i >= j {
			break
		}

		swap(arr, i, j)
	}

	// 最后将pivot的值，放到一个左边元素比它小，右边元素比它大的地方
	swap(arr, lo, j)
	return j
}

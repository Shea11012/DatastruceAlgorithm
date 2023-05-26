package sort

import (
	"math/rand"
	"time"
)

func swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func shuffle(arr []int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < len(arr); i++ {
		r := i + r.Intn(len(arr)-i)
		swap(arr, i, r)
	}
}

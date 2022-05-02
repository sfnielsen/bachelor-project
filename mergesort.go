package main

import (
	"sync"
)

func merge(a []int, b []int) []int {
    final := []int{}
    i := 0
    j := 0
    for i < len(a) && j < len(b) {
        if a[i] < b[j] {
            final = append(final, a[i])
            i++
        } else {
            final = append(final, b[j])
            j++
        }
    }
    for ; i < len(a); i++ {
        final = append(final, a[i])
    }
    for ; j < len(b); j++ {
        final = append(final, b[j])
    }
    return final
}

func mergeSort(items []int) []int {
    if len(items) < 2 {
        return items
    }
    first := mergeSort(items[:len(items)/2])
    second := mergeSort(items[len(items)/2:])
    return merge(first, second)
}

func mergesortv3(s []int) {
	len := len(s)

	if len > 1 {
		if len <= 2048 { // Sequential
			mergeSort(s)
		} else { // Parallel
			middle := len / 2

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				mergesortv3(s[:middle])
			}()

			mergesortv3(s[middle:])

			wg.Wait()
			merge(s, middle)
		}
	}
}
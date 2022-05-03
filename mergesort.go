package main

import (
	"sort"
)

var partition int

func merge(a []Tuple, b []Tuple) []Tuple {
	final := []Tuple{}
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i].value < b[j].value {
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

func MergeSort(data []Tuple, r chan []Tuple) {
	if len(data) == 1 {
		r <- data
		return
	}
	if len(data) <= 2048 { // Sequential
		mergeSortSeq(data)
	}

	leftChan := make(chan []Tuple)
	rightChan := make(chan []Tuple)
	middle := len(data) / 2

	go MergeSort(data[:middle], leftChan)
	go MergeSort(data[middle:], rightChan)

	ldata := <-leftChan
	rdata := <-rightChan

	close(leftChan)
	close(rightChan)
	r <- merge(ldata, rdata)
	return
}

func ParallelQuick(data []Tuple, r chan []Tuple) {
	if len(data) == 1 {
		r <- data
		return
	}
	if len(data) <= partition { // Sequential
		sort.Slice(data, func(a, b int) bool {
			return (data[a].value < data[b].value)
		})
		r <- data
		return
	}
	leftChan := make(chan []Tuple)
	rightChan := make(chan []Tuple)
	middle := len(data) / 2

	go ParallelQuick(data[:middle], leftChan)
	go ParallelQuick(data[middle:], rightChan)

	ldata := <-leftChan
	rdata := <-rightChan

	close(leftChan)
	close(rightChan)
	r <- merge(ldata, rdata)
	return
}

func Quick(data []Tuple, c chan []Tuple) {
	sort.Slice(data, func(a, b int) bool {
		return (data[a].value < data[b].value)
	})
	c <- data
}

/*
func mergeSortOnceParallel(s []Tuple, r chan []Tuple) []Tuple {
	lenghtOfS := len(s)
	s1 := make(chan []Tuple)

	middle := lenghtOfS / 2

	var wg sync.WaitGroup
	wg.Add(1)

	leftChan := make(chan []Tuple)
	rightChan := make(chan []Tuple)

	go MergeSort(s[middle:], leftChan)
	fmt.Println("first send")
	go MergeSort(s[:middle], rightChan)
	fmt.Println("second send")

	wg.Wait()

	fmt.Println("done")
	s1_res := <-s1

	return merge(s1_res, s2)
}

func mergeSortParallel(s []Tuple) []Tuple {
	s1 := make(chan []Tuple)

	lenghtOfS := len(s)
	if lenghtOfS > 1 {
		if lenghtOfS <= 2048 { // Sequential
			return MergeSort(s)
		} else { // Parallel
			middle := lenghtOfS / 2

			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				defer wg.Done()
				t := s
				t = t[:middle]
				mergeSortParallel(t)
			}()
			t2 := s
			t2 = t2[:middle]
			s2 := mergeSortParallel(t2)

			wg.Wait()

			s11 := <-s1
			return merge(s11, s2)
		}
	}
	return s
}
*/
func mergeSortSeq(items []Tuple) []Tuple {
	if len(items) < 2 {
		return items
	}
	first := mergeSortSeq(items[:len(items)/2])
	second := mergeSortSeq(items[len(items)/2:])
	return mergeSeq(first, second)
}

func mergeSeq(a []Tuple, b []Tuple) []Tuple {
	final := []Tuple{}
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i].value < b[j].value {
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

func mergesortPara(s []Tuple) {
	len := len(s)
	res := make(chan []Tuple)
	if len > 1 {
		if len <= 2048 { // Sequential
			mergeSortSeq(s)
		} else { // Parallel
			go MergeSort(s, res)
		}
	}
}

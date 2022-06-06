package main

import (
	"bytes"
	"encoding/binary"
)

const digit = 8
const maxbit = -1 << 31

/*
func main() {
	var data = []float64{421.123, 15.121231233, 175.123123123, 90.123, 2.11, 214.33, 52.11, 166.33, 123123123.3333333}
	fmt.Println("\n--- Unsorted --- \n\n", data)
	radixsort(data)
	fmt.Println("\n--- Sorted ---\n\n", data, "\n")
}
*/

func radixsort(data []Tuple) {
	buf := bytes.NewBuffer(nil)
	ds := make([][]byte, len(data))

	for i, e := range data {
		binary.Write(buf, binary.LittleEndian, e.value)
		b := make([]byte, digit)
		buf.Read(b)
		ds[i] = b
	}

	countingSort := make([][][]byte, 256)
	tuples := make([][]Tuple, 256)

	for i := 0; i < digit; i++ {
		for asdf, b := range ds {
			countingSort[b[i]] = append(countingSort[b[i]], b)
			tuples[b[i]] = append(tuples[b[i]], data[asdf])
		}
		j := 0
		for k, bs := range countingSort {
			copy(ds[j:], bs)
			copy(data[j:], tuples[k])

			j += len(bs)

			countingSort[k] = bs[:0]
			tuples[k] = tuples[k][:0]

		}
	}
	/*
		var w float64

		for i, b := range ds {
			buf.Write(b)
			binary.Read(buf, binary.LittleEndian, &w)
			data[i] = placementMap[w]
		}
	*/
}

package main

import (
	"fmt"
	"math"
)

func main() {
	D := [][]float64{
		{0, 17, 21, 27},
		{17, 0, 12, 18},
		{21, 12, 0, 14},
		{27, 18, 14, 0},
	}

	neighbourJoin(D)
}

func neighbourJoin(D [][]float64) {
	n := len(D)
	M := make([][]float64, n)
	for i := range M {
		M[i] = make([]float64, n)
	}
	r := make([]float64, n)
	print(len(M[0]))

	print("D\n")
	for i := 0; i < n; i++ {

		fmt.Println(D[i])
	}

	for i, row := range M {
		sum := 0.0
		for j := range row {
			sum = sum + D[i][j]
			print(D[i][j])
			print("\n")
		}
		r[i] = sum / float64(n-2)
		print(r[i])
		print("\n")
	}

	cur_val := math.MaxFloat64
	cur_i, cur_j := -1, -1

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {

			if i == j {
				M[i][j] = 0
			} else {
				M[i][j] = D[i][j] - r[i] - r[j]

				if M[i][j] < cur_val {
					cur_val = M[i][j]
					cur_i = i
					cur_j = j
				}
			}
		}
	}

	v_iu := D[cur_i][cur_j]/2 + (r[cur_i]-r[cur_j])/2
	v_ju := D[cur_i][cur_j]/2 + (r[cur_j]-r[cur_i])/2

	print("\n")
	print(v_iu, v_ju)

	print("M\n")

	for i := 0; i < n; i++ {
		fmt.Println(M[i])
	}

}

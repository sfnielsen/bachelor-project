package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
)

var NewickFlag bool = true

func main() {
	D := [][]float64{
		{0, 17, 21, 27},
		{17, 0, 12, 18},
		{21, 12, 0, 14},
		{27, 18, 14, 0},
	}
	n := len(D)
	Q := make([][]float64, n)
	for i := range Q {
		Q[i] = make([]float64, n)
	}

	labels := []string{
		"A", "B", "C", "D",
	}
	neighborJoin(Q, D, labels)
}

func canonicalNeighborJoining(M [][]float64, r []float64, D [][]float64, n int) (int, int) {
	cur_val := -math.MaxFloat64
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
	return cur_i, cur_j

}

type Tuple struct {
	value   float64
	index_j int
}

func MaxIntSlice(v []float64) (m float64) {
	m = -math.MaxFloat64

	for i := 0; i < len(v); i++ {
		if v[i] > m {
			m = v[i]
		}
	}
	return m
}

func rapidNeighborJoining(u []float64, D [][]float64, n int) (int, int) {
	fmt.Println("swampgod")
	S := make([][]Tuple, n)
	for i := range S {
		S[i] = make([]Tuple, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var tuple Tuple
			tuple.value = D[i][j]
			tuple.index_j = j

			S[i][j] = tuple

		}

		//sorting row in S
		sort.Slice(S[i], func(a, b int) bool {
			return (S[i][a].value < S[i][b].value)
		})

		fmt.Println(S[i])
		fmt.Println("sawwaw")
	}

	max_u := MaxIntSlice(u)
	q_min := math.MaxFloat64
	cur_i, cur_j := -1, -1

	for r := 0; r < n; r++ {
		for c := 1; c < n; c++ {
			s := S[r][c]
			if s.value-u[r]-max_u > q_min {
				break
			}
			if s.value-u[r]-u[c] < q_min {
				cur_i = r
				cur_j = s.index_j
				q_min = s.value - u[r] - u[c]
			}
		}
	}

	return cur_i, cur_j
}

func neighborJoin(M [][]float64, D [][]float64, labels []string) {
	n := len(D)

	u := make([]float64, n)

	print("D\n")
	for i := 0; i < n; i++ {
		fmt.Println(D[i])
	}
	print("\n")
	for i, row := range M {
		sum := 0.0
		for j := range row {
			sum = sum + D[i][j]
		}
		u[i] = sum / float64(n-2)
	}

	cur_i, cur_j := rapidNeighborJoining(u, D, n)

	if NewickFlag {
		//Distance to new point where they meet
		v_iu := fmt.Sprintf("%f", D[cur_i][cur_j]/2+(u[cur_i]-u[cur_j])/2)
		v_ju := fmt.Sprintf("%f", D[cur_i][cur_j]/2+(u[cur_j]-u[cur_i])/2)
		//convert to string
		fmt.Println(v_iu)
		fmt.Println(v_iu)

		//creating newick form
		labels[cur_i] = "(" + labels[cur_i] + ":" + v_iu + "," + labels[cur_j] + ":" + v_ju + ")"
		labels = append(labels[:cur_j], labels[cur_j+1:]...)
	}

	D_new := createNewDistanceMatrix(D, cur_i, cur_j)

	for i := 0; i < len(labels); i++ {
		fmt.Println(labels[i])
	}

	//stop maybe
	if len(D_new) > 2 {
		neighborJoin(M, D_new, labels)
	} else {
		if NewickFlag {
			newick := "(" + labels[cur_i] + ":" + fmt.Sprintf("%f", D_new[cur_i][cur_j]/2) + "," + labels[cur_j] + ":" + fmt.Sprintf("%f", D_new[cur_i][cur_j]/2) + ");"
			fmt.Println(newick)

			err := ioutil.WriteFile("newick.txt", []byte(newick), 0644)
			if err != nil {
				panic(err)
			}

		}
		return
	}

}

func createNewDistanceMatrix(D [][]float64, p_i int, p_j int) [][]float64 {

	for k := 0; k < len(D); k++ {
		if p_i == k {
			D[p_i][k] = 0
		}
		if p_j == k {
			continue
		} else {
			//Overwrite p_i as merge ij
			temp := (D[p_i][k] + D[p_j][k] - D[p_i][p_j]) / 2
			D[p_i][k] = temp
			D[k][p_i] = temp
		}
	}

	D_new := append(D[:p_j], D[p_j+1:]...)

	for i := 0; i < len(D_new); i++ {
		D_new[i] = append(D_new[i][:p_j], D_new[i][p_j+1:]...)
	}

	return D_new
}

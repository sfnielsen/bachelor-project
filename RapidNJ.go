package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
)

var NewickFlag bool = true

func main() {
	labels := []string{
		"A", "B", "C", "D",
	}
	D := [][]float64{
		{0, 17, 21, 27},
		{17, 0, 12, 18},
		{21, 12, 0, 14},
		{27, 18, 14, 0},
	}
	n := len(D)

	//create initial S matrix
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
	}
	dead_records := make(map[int]bool)

	neighborJoin(D, S, labels, dead_records)
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

func rapidNeighborJoining(u []float64, D [][]float64, S [][]Tuple) (int, int) {
	n := len(D)
	fmt.Println("swampgod")

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

func neighborJoin(D [][]float64, S [][]Tuple, labels []string, dead_records map[int]bool) {
	n := len(D)

	u := make([]float64, n)

	print("D\n")
	for i := 0; i < n; i++ {
		fmt.Println(D[i])
	}

	print("\n")
	for i, row := range D {
		sum := 0.0
		for j := range row {
			sum = sum + D[i][j]
		}
		u[i] = sum / float64(n-2)
	}

	cur_i, cur_j := rapidNeighborJoining(u, D, S)

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

	D_new, S_new, dead_records_new := createNewDistanceMatrix(S, dead_records, D, cur_i, cur_j)

	for i := 0; i < len(labels); i++ {
		fmt.Println(labels[i])
	}

	//stop maybe
	if len(D_new) > 2 {
		neighborJoin(D_new, S_new, labels, dead_records_new)
	} else {
		if NewickFlag {
			fmt.Println(cur_i, cur_j)
			newick := "(" + labels[0] + ":" + fmt.Sprintf("%f", D_new[0][1]/2) + "," + labels[1] + ":" + fmt.Sprintf("%f", D_new[0][1]/2) + ");"
			fmt.Println(newick)

			err := ioutil.WriteFile("newick.txt", []byte(newick), 0644)
			if err != nil {
				panic(err)
			}

		}
		return
	}

}

func createNewDistanceMatrix(S [][]Tuple, dead_records map[int]bool, D [][]float64, p_i int, p_j int) ([][]float64, [][]Tuple, map[int]bool) {
	for k := 0; k < len(D); k++ {
		if p_i == k {
			continue
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

	//delete row in both D and S
	D_new := append(D[:p_j], D[p_j+1:]...)

	//delete column in D
	for i := 0; i < len(D_new); i++ {
		D_new[i] = append(D_new[i][:p_j], D_new[i][p_j+1:]...)

	}

	//fix S
	S_new := S

	//overwrite the row p_i where we want to store merged ij
	fmt.Println(p_i, p_j, len(D[p_i]))
	fmt.Println("HAHA")
	for i := 0; i < len(S_new); i++ {
		fmt.Println(S_new[i])
	}
	for j := 0; j < len(D[p_i]); j++ {
		var tuple Tuple
		tuple.value = D[p_i][j]
		tuple.index_j = j
		S_new[p_i][j] = tuple

	}
	//cut excess data away
	S_new[p_i] = S_new[p_i][:len(D)-1]

	//sort merged row
	sort.Slice(S_new[p_i], func(a, b int) bool {
		return (S_new[p_i][a].value < S_new[p_i][b].value)
	})

	fmt.Println("asdhf")
	for i := 0; i < len(S_new); i++ {
		fmt.Println(S_new[i])
	}
	S_new = append(S[:p_j], S[p_j+1:]...)
	dead_records[p_j] = true

	return D_new, S_new, dead_records
}

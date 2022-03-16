package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

func canonicalNeighborJoining(Q [][]float64, r []float64, D [][]float64, n int) (int, int) {
	cur_val := math.MaxFloat64
	cur_i, cur_j := -1, -1

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {

			if i == j {
				Q[i][j] = 0
			} else {
				Q[i][j] = D[i][j] - r[i] - r[j]

				if Q[i][j] < cur_val {
					cur_val = Q[i][j]
					cur_i = i
					cur_j = j
				}
			}
		}
	}
	return cur_i, cur_j

}

func neighborJoin(D [][]float64, labels []string, array Tree, tree Tree) (string, Tree) {
	n := len(D)
	Q := make([][]float64, n)
	for i := range Q {
		Q[i] = make([]float64, n)
	}
	u := make([]float64, n)

	print("D\n")
	for i := 0; i < n; i++ {
		fmt.Println(D[i])
	}
	print("\n")
	for i, row := range Q {
		sum := 0.0
		for j := range row {
			sum = sum + D[i][j]
		}
		u[i] = sum / float64(n-2)
	}

	cur_i, cur_j := canonicalNeighborJoining(Q, u, D, n)

	if NewickFlag {
		//Distance to new point where they meet
		v_iu := fmt.Sprintf("%f", D[cur_i][cur_j]/2+(u[cur_i]-u[cur_j])/2)
		v_ju := fmt.Sprintf("%f", D[cur_i][cur_j]/2+(u[cur_j]-u[cur_i])/2)
		//convert to string
		fmt.Println(v_iu)
		fmt.Println(v_iu)

		distance_to_x, _ := strconv.ParseFloat(v_iu, 64)
		distance_to_y, _ := strconv.ParseFloat(v_ju, 64)

		newNode := integrateNewNode(array[cur_i], array[cur_j], distance_to_x, distance_to_y)
		array[cur_i] = newNode
		tree = append(tree, newNode)

		array = append(array[:cur_j], array[cur_j+1:]...)

		//creating newick form
		labels[cur_i] = "(" + labels[cur_i] + ":" + v_iu + "," + labels[cur_j] + ":" + v_ju + ")"
		labels = append(labels[:cur_j], labels[cur_j+1:]...)
	}

	D_new := createNewDistanceMatrixNJ(D, cur_i, cur_j)

	for i := 0; i < len(labels); i++ {
		fmt.Println(labels[i])
	}

	//stop maybe
	if len(D_new) > 2 {
		neighborJoin(D_new, labels, array, tree)
	} else {
		if NewickFlag {
			fmt.Println(cur_i, cur_j)
			newick := "(" + labels[0] + ":" + fmt.Sprintf("%f", D_new[0][1]/2) + "," + labels[1] + ":" + fmt.Sprintf("%f", D_new[0][1]/2) + ");"
			fmt.Println(newick)

			err := ioutil.WriteFile("newick.txt", []byte(newick), 0644)
			if err != nil {
				panic(err)
			}

			new_edge_0 := new(Edge)
			new_edge_0.Distance = D_new[0][1]
			new_edge_0.Node = array[1]

			new_edge_1 := new(Edge)
			new_edge_1.Distance = D_new[0][1]
			new_edge_1.Node = array[0]

			array[0].Edge_array = append(array[0].Edge_array, new_edge_0)
			array[1].Edge_array = append(array[1].Edge_array, new_edge_1)

			array = remove(array, 0)

			return newick, tree
		}
	}
	return "error", tree
}

func createNewDistanceMatrixNJ(D [][]float64, p_i int, p_j int) [][]float64 {

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

	D_new := append(D[:p_j], D[p_j+1:]...)

	for i := 0; i < len(D_new); i++ {
		D_new[i] = append(D_new[i][:p_j], D_new[i][p_j+1:]...)
	}

	return D_new
}

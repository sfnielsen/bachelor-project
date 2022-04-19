package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
)

var NewickFlag bool = true

//function to initialize dead records
func initDeadRecords(D [][]float64) map[int]int {
	live_records := make(map[int]int)
	for i := range D {
		live_records[i] = i
	}
	return live_records
}

//function to initialize S matrix
func initSmatrix(D [][]float64) [][]Tuple {
	n := len(D)
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
	}
	return S
}

type Tuple struct {
	value   float64
	index_j int
}

type U_Tuple struct {
	value      float64
	index_in_d int
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

func generateTreeForRapidNJ(labels []string) Tree {
	tree := make(Tree, 0)

	for _, label := range labels {
		node_to_append := new(Node)
		node_to_append.Name = label
		tree = append(tree, node_to_append)
	}
	return tree
}

//function to create initial u array
func create_u(D [][]float64) []float64 {
	n := len(D)
	u := make([]float64, n)
	for i, row := range D {
		sum := 0.0
		for j := range row {
			sum += D[i][j]
		}
		u[i] = float64(sum) / float64(n-2)
	}

	return u
}

//function to update u instead of initializing it each iteration.
//this results in lacking precision due to float pointer errors
//but speeds up the running time significantly
func update_u(D [][]float64, u []float64, i int, j int) []float64 {

	n := len(D)

	//update i with merge ij
	u[i] = 0

	for idx := range u {
		if idx == i || idx == j {
			continue
		}
		u[idx] = u[idx]*float64(n-2) - D[idx][i] - D[idx][j]
		new_dist := (D[i][idx] + D[j][idx] - D[i][j]) / 2.0

		u[idx] += new_dist
		u[idx] /= float64(n - 3)

		//also add value to ij merge
		u[i] += new_dist

	}
	u[i] /= float64(n - 3)

	//remove j from the array
	u = append(u[:j], u[j+1:]...)

	return u
}

func rapidNeighbourJoin(D [][]float64, labels []string, s_search_strategy S_Search_Strategy) (string, Tree) {

	//setup initial data structures for rapidNJ
	S := initSmatrix(D)
	deadRecords := initDeadRecords(D)
	var label_tree Tree = generateTreeForRapidNJ(labels)
	var tree Tree
	tree = append(tree, label_tree...)
	total_nodes := len(S) - 1

	u := create_u(D)

	//run rapidNJ algorithm
	newick, tree := rapidJoinRec(D, S, labels, deadRecords, label_tree, tree, total_nodes, u, s_search_strategy)

	return newick, tree

}

//two Tree types. array Tree manages connection between labels and matrix while tree Tree holds all nodes (tips AND INTERNALS)
func rapidJoinRec(D [][]float64, S [][]Tuple, labels []string, live_records map[int]int, array Tree, tree Tree, total_nodes int, u []float64,
	s_search_strategy S_Search_Strategy) (string, Tree) {

	//gets two indexes in D
	cur_i, cur_j := s_search_strategy(u, D, S, live_records)

	//make sure cur_i is the smallest index.
	//both important for labels and for creation of new distance matrix.
	if cur_i > cur_j {
		temp := cur_i
		cur_i = cur_j
		cur_j = temp
	}

	if NewickFlag {
		//Distance to new point where they meet
		v_iu := fmt.Sprintf("%f", D[cur_i][cur_j]/2+(u[cur_i]-u[cur_j])/2)
		v_ju := fmt.Sprintf("%f", D[cur_i][cur_j]/2+(u[cur_j]-u[cur_i])/2)
		//convert to string

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

	//update u for next iteration
	u = update_u(D, u, cur_i, cur_j)

	//update total nodes for next iteration
	total_nodes++

	D_new, S_new, live_records_new := createNewDistanceMatrix(S, live_records, D, cur_i, cur_j, total_nodes)

	//stop maybe
	if len(D_new) > 2 {
		return rapidJoinRec(D_new, S_new, labels, live_records_new, array, tree, total_nodes, u, s_search_strategy)
	} else {
		if NewickFlag {
			newick := "(" + labels[0] + ":" + fmt.Sprintf("%f", D_new[0][1]/2) + "," + labels[1] + ":" + fmt.Sprintf("%f", D_new[0][1]/2) + ");"
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
	return "error", tree //this case should not be possible
}

func createNewDistanceMatrix(S [][]Tuple, live_records map[int]int, D [][]float64, p_i int, p_j int, new_map_key int) ([][]float64, [][]Tuple, map[int]int) {

	for k := 0; k < len(D); k++ {
		if p_i == k {
			continue
		}
		if p_j == k {
			continue
		} else {
			//Overwrite p_i as merge ij
			temp := (D[p_i][k] + D[p_j][k] - D[p_i][p_j]) / 2.0
			D[p_i][k] = temp
			D[k][p_i] = temp

		}
	}

	//delete row in D
	D_new := append(D[:p_j], D[p_j+1:]...)

	//delete column in D
	for i := 0; i < len(D_new); i++ {
		D_new[i] = append(D_new[i][:p_j], D_new[i][p_j+1:]...)
	}

	//Overwrite dead records
	for k, v := range live_records {
		if _, ok := live_records[k]; !ok {
			continue
		}
		if v == p_i {
			delete(live_records, k)
		} else if v == p_j {
			delete(live_records, k)
		} else if v > p_j {
			live_records[k] = v - 1
		}

	}
	live_records[new_map_key] = p_i

	//allow quicker lookups
	live_records_reverse := reverseMap((live_records))

	//fix S
	S_new := S

	//overwrite the row p_i where we want to store merged ij
	for j := 0; j < len(D[p_i]); j++ {
		var tuple Tuple
		var result int

		tuple.value = D[p_i][j]
		result = live_records_reverse[j]

		tuple.index_j = result
		S_new[p_i][j] = tuple

	}
	//cut excess data away
	S_new[p_i] = S_new[p_i][:len(D_new)]

	//sort merged row
	sort.Slice(S_new[p_i], func(a, b int) bool {
		return (S_new[p_i][a].value < S_new[p_i][b].value)
	})

	//delete row in S
	S_new = append(S[:p_j], S[p_j+1:]...)

	return D_new, S_new, live_records
}

func reverseMap(m map[int]int) map[int]int {
	n := make(map[int]int, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

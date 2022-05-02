package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"sync"
)

var NewickFlag bool = true

//function to initialize dead records
func initLiveRecords(D [][]float64) map[int]int {
	live_records := make(map[int]int)
	for i := range D {
		live_records[i] = i
	}
	return live_records
}

//function to initialize S matrix
func initSmatrix(D [][]float64) [][]Tuple {
	var wg sync.WaitGroup

	n := len(D)
	S := make([][]Tuple, n)
	wg.Add(n)

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
		go_i := i
		go sort_S_row(&wg, &S[go_i])

	}

	wg.Wait()

	return S
}

func sort_S_row(wg *sync.WaitGroup, row *[]Tuple) {
	defer wg.Done()
	row_sort := *row
	sort.Slice(row_sort, func(a, b int) bool {
		return (row_sort[a].value < row_sort[b].value)
	})

	*row = row_sort

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
//This speeds up the running time significantly
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
	u[j] = 0
	u[j] = u[len(u)-1]
	u[len(u)-1] = 0
	u = u[:len(u)-1]
	//u = append(u[:j], u[j+1:]...)

	return u
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

func rapidNeighbourJoin(D [][]float64, labels []string, s_search_strategy S_Search_Strategy) (string, Tree) {

	//setup initial data structures for rapidNJ
	S := initSmatrix(D)

	liveRecords := initLiveRecords(D)
	liveRecordsReverse := reverseMap(liveRecords)

	var label_tree Tree = generateTreeForRapidNJ(labels)
	var tree Tree
	tree = append(tree, label_tree...)
	total_nodes := len(S) - 1

	u := create_u(D)

	//run rapidNJ algorithm
	newick, tree := rapidJoinRec(D, S, labels, liveRecords, liveRecordsReverse, label_tree, tree, total_nodes, u, s_search_strategy, -1)

	return newick, tree

}

//two Tree types. array Tree manages connection between labels and matrix while tree Tree holds all nodes (tips AND INTERNALS)
func rapidJoinRec(D [][]float64, S [][]Tuple, labels []string, live_records map[int]int, live_records_reverse map[int]int,
	array Tree, tree Tree, total_nodes int, u []float64, s_search_strategy S_Search_Strategy, row_for_next_it int) (string, Tree) {

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

		//array = append(array[:cur_j], array[cur_j+1:]...)

		array[cur_j] = array[len(array)-1]
		array[len(array)-1] = nil
		array = array[:len(array)-1]

		//creating newick form
		labels[cur_i] = "(" + labels[cur_i] + ":" + v_iu + "," + labels[cur_j] + ":" + v_ju + ")"

		//labels = append(labels[:cur_j], labels[cur_j+1:]...)

		labels[cur_j] = labels[len(labels)-1]
		labels[len(labels)-1] = ""
		labels = labels[:len(labels)-1]

	}

	//update u for next iteration
	u = update_u(D, u, cur_i, cur_j)

	//update total nodes for next iteration
	total_nodes++

	D_new, S_new, live_records_new, live_records_reverse_new := createNewDistanceMatrix(S, live_records, live_records_reverse, D, cur_i, cur_j, total_nodes)

	//stop maybe
	if len(D_new) > 2 {
		return rapidJoinRec(D_new, S_new, labels, live_records_new, live_records_reverse_new, array, tree, total_nodes, u, s_search_strategy, row_for_next_it)
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

func createNewDistanceMatrix(S [][]Tuple, live_records map[int]int, live_records_reverse map[int]int,
	D [][]float64, p_i int, p_j int, new_map_key int) ([][]float64, [][]Tuple, map[int]int, map[int]int) {

	//update D
	D_new := update_D(D, p_i, p_j)

	//Overwrite live records
	live_records = update_live_records(live_records, p_i, p_j, new_map_key)

	//allow quicker lookups when updating S
	live_records_reverse = update_live_records_reverse(live_records_reverse, p_i, p_j, new_map_key)

	//update S
	S_new := update_S(S, D, p_i, p_j, live_records_reverse)

	return D_new, S_new, live_records, live_records_reverse
}

func update_S(S [][]Tuple, D [][]float64, p_i int, p_j int, live_records_reverse map[int]int) [][]Tuple {
	S_new := S
	//overwrite the row p_i where we want to store merged ij
	s_sorting_indexes := make([]string, len(S))
	for j := 0; j < len(D[p_i]); j++ {
		var tuple Tuple
		var result int

		tuple.value = D[p_i][j]
		result = live_records_reverse[j]
		s_sorting_indexes[j] = convertFloatToIntString(D[p_i][j])

		tuple.index_j = result
		S_new[p_i][j] = tuple

	}
	//cut excess data away
	S_new[p_i] = S_new[p_i][:len(D)-1]

	//sort merged
	Sort(s_sorting_indexes)
	fmt.Println(s_sorting_indexes)
	sort.Slice(S_new[p_i], func(a, b int) bool {
		return (S_new[p_i][a].value < S_new[p_i][b].value)
	})

	//delete row in S
	//S_new = append(S[:p_j], S[p_j+1:]...)

	S_new[p_j] = nil
	S_new[p_j] = S_new[len(S_new)-1]
	S_new[p_j] = S_new[p_j][:len(S_new[len(S_new)-1])]
	S_new[len(S_new)-1] = nil
	S_new = S_new[:len(S_new)-1]

	return S_new
}

func update_D(D [][]float64, p_i int, p_j int) [][]float64 {

	//update i rows/cols
	D = update_row_col_i(D, p_i, p_j)

	//remove j rows/cols
	D = update_row_col_j(D, p_j)

	return D
}

func update_row_col_i(D [][]float64, p_i int, p_j int) [][]float64 {

	row_to_delete := make([]float64, len(D[0]))
	copy(row_to_delete, D[p_j])

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
	return D
}

func update_row_col_j(D [][]float64, p_j int) [][]float64 {

	//delete last row in D and move it to p_j
	D[p_j] = D[len(D)-1]
	D[len(D)-1] = nil
	D = D[:len(D)-1]

	//delete last column in D and move it to p_j
	for i := 0; i < len(D); i++ {
		D[i][p_j] = D[i][len(D[i])-1]
		D[i][len(D[i])-1] = 0
		D[i] = D[i][:len(D[i])-1]
	}

	return D
}

func update_live_records(live_records map[int]int, p_i int, p_j int, new_map_key int) map[int]int {

	var last_k, last_v int

	last_v = len(live_records) - 1

	for k, v := range live_records {
		if _, ok := live_records[k]; !ok {
			continue
		}
		if v == last_v {
			last_k = k
		}
		if v == p_j {
			delete(live_records, k)
		}
		if v == p_i {
			delete(live_records, k)
		}

		//else if v > p_j {
		//	live_records[k] = v - 1

	}

	if _, ok := live_records[last_k]; ok {
		live_records[last_k] = p_j
	}
	live_records[new_map_key] = p_i

	return live_records
}

func update_live_records_reverse(live_records_reverse map[int]int, p_i int, p_j int, new_map_key int) map[int]int {

	//for x := p_j + 1; x < len(live_records_reverse); x++ {
	//	live_records_reverse[x-1] = live_records_reverse[x]
	//}

	live_records_reverse[p_i] = new_map_key
	live_records_reverse[p_j] = live_records_reverse[len(live_records_reverse)-1]
	delete(live_records_reverse, len(live_records_reverse)-1)
	//delete(live_records_reverse, len(live_records_reverse)-1)

	return live_records_reverse
}

func reverseMap(m map[int]int) map[int]int {
	n := make(map[int]int, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

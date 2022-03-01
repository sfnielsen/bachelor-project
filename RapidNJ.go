package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
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

	S := initSmatrix(D)
	dead_records := initDeadRecords(D)

	var treeBanana Tree
	var array Tree
	if NewickFlag {
		array = generateTreeForRapidNJ(labels)
		for _, node := range array {
			treeBanana = append(treeBanana, node)
		}
	} else {
		array = make(Tree, 0)
	}

	newick_result, _ := neighborJoin(D, S, labels, dead_records, array, treeBanana)
	fmt.Println(newick_result)
}

//function to initialize dead records
func initDeadRecords(D [][]float64) map[int]int {
	dead_records := make(map[int]int)
	for i := range D {
		dead_records[i] = i
	}
	return dead_records
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
		fmt.Println(S[i])
	}
	return S
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

func rapidNeighborJoining(u []float64, D [][]float64, S [][]Tuple, dead_records map[int]int) (int, int) {
	max_u := MaxIntSlice(u)
	q_min := math.MaxFloat64
	cur_i, cur_j := -1, -1

	fmt.Println("darecords", dead_records)

	for r, row := range S {
		fmt.Println("hey", row, len(row))

		for i, _ := range row {
			if i == 0 {
				fmt.Println("I==0 case")
				continue
			}
			s := S[r][i]

			c_to_cD := dead_records[s.index_j]

			//check if dead record
			if c_to_cD == -1 {
				fmt.Println("C_TO_CD CASE")
				continue
			}
			// case where i == j
			if r == c_to_cD {
				fmt.Println("I==J CASE")
				continue
			}
			if s.value-u[r]-max_u > q_min {
				break
			}
			if s.value-u[r]-u[c_to_cD] < q_min {
				cur_i = r
				cur_j = s.index_j
				q_min = s.value - u[r] - u[c_to_cD]
			}
		}
	}

	return cur_i, cur_j
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

func neighborJoin(D [][]float64, S [][]Tuple, labels []string, dead_records map[int]int, array Tree, treeBanana Tree) (string, Tree) {

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

	cur_i, cur_j := rapidNeighborJoining(u, D, S, dead_records)

	j_in_D := dead_records[cur_j]

	fmt.Println("indices", cur_i, cur_j)

	if NewickFlag {

		if cur_i == -1 || j_in_D == -1 {
			fmt.Println(cur_i, j_in_D, "BABABBBA tis", dead_records, len(D))
		}
		//Distance to new point where they meet
		v_iu := fmt.Sprintf("%f", D[cur_i][j_in_D]/2+(u[cur_i]-u[j_in_D])/2)
		v_ju := fmt.Sprintf("%f", D[cur_i][j_in_D]/2+(u[cur_j]-u[j_in_D])/2)
		//convert to string
		fmt.Println(v_iu)
		fmt.Println(v_iu)

		//make sure p_i is the smallest index and dont change it w.r.t newick implementation
		temp_i := 0
		temp_j := 0
		if cur_i > j_in_D {
			temp_i = j_in_D
			temp_j = cur_i
		} else {
			temp_i = cur_i
			temp_j = j_in_D
		}

		distance_to_y, _ := strconv.ParseFloat(v_iu, 64)
		distance_to_x, _ := strconv.ParseFloat(v_ju, 64)

		newNode := integrateNewNode(array[temp_i], array[temp_j], distance_to_x, distance_to_y)
		array[temp_i] = newNode
		treeBanana = append(treeBanana, newNode)
		array = append(array[:temp_j], array[temp_j+1:]...)

		//creating newick form
		labels[cur_i] = "(" + labels[cur_i] + ":" + v_iu + "," + labels[j_in_D] + ":" + v_ju + ")"
		labels = append(labels[:j_in_D], labels[j_in_D+1:]...)

		for i, v := range labels {
			fmt.Println(i, v, "This is god")
		}
	}

	D_new, S_new, dead_records_new := createNewDistanceMatrix(S, dead_records, D, cur_i, cur_j)
	fmt.Println("newrecords", dead_records_new)
	for i := 0; i < len(labels); i++ {
		fmt.Println(labels[i])
	}

	//stop maybe
	if len(D_new) > 2 {
		return neighborJoin(D_new, S_new, labels, dead_records_new, array, treeBanana)
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

			return newick, treeBanana

		}
	}
	return "error", treeBanana //this case should not be possible
}

func createNewDistanceMatrix(S [][]Tuple, dead_records map[int]int, D [][]float64, p_i int, p_j int) ([][]float64, [][]Tuple, map[int]int) {
	//make sure p_i is the smallest index

	p_j_in_D := dead_records[p_j]

	if p_i > p_j_in_D {
		temp := p_i
		p_i = p_j_in_D
		p_j_in_D = temp
	}

	for k := 0; k < len(D); k++ {
		if p_i == k {
			continue
		}
		if p_j_in_D == k {
			continue
		} else {
			//Overwrite p_i as merge ij
			temp := (D[p_i][k] + D[p_j_in_D][k] - D[p_i][p_j_in_D]) / 2
			D[p_i][k] = temp
			D[k][p_i] = temp

		}
	}

	//delete row in both D and S
	D_new := append(D[:p_j_in_D], D[p_j_in_D+1:]...)

	//delete column in D
	for i := 0; i < len(D_new); i++ {
		D_new[i] = append(D_new[i][:p_j_in_D], D_new[i][p_j_in_D+1:]...)

	}

	//fix S
	S_new := S

	//overwrite the row p_i where we want to store merged ij
	for j := 0; j < len(D[p_i]); j++ {
		var tuple Tuple
		tuple.value = D[p_i][j]
		tuple.index_j = j
		S_new[p_i][j] = tuple

	}
	//cut excess data away
	S_new[p_i] = S_new[p_i][:len(D)]

	//sort merged row
	sort.Slice(S_new[p_i], func(a, b int) bool {
		return (S_new[p_i][a].value < S_new[p_i][b].value)
	})

	S_new = append(S[:p_j], S[p_j+1:]...)

	//assign dead records -> -1
	dead_records[p_i] = -1
	dead_records[p_j] = -1
	//add merged ij at i's spot
	for k, v := range dead_records {
		//if affected by index movement
		if k > p_j {
			//if record already dead we keep -1 as the 'nil' value
			if dead_records[k] != -1 {
				dead_records[k] = v - 1
			}
		}
	}
	dead_records[len(dead_records)] = p_i

	return D_new, S_new, dead_records
}

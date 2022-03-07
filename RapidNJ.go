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

	for r, row := range S {
		fmt.Println("hey", row, len(row))

		for c := range row {
			s := S[r][c]
			c_to_cD := dead_records[s.index_j]
			fmt.Println(c, s.value, c_to_cD)
			//check if dead record
			if c_to_cD == -1 {
				continue
			}
			// case where i == j
			if r == c_to_cD {
				continue
			}
			if s.value-u[r]-max_u > q_min {
				fmt.Println("breaking")
				break
			}
			fmt.Println("some length should be the same; ", len(u), len(D), len(S))
			q := s.value - u[r] - u[c_to_cD]
			if q < q_min {
				cur_i = r
				cur_j = c_to_cD
				q_min = q
				fmt.Println(q_min, "lautis")
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

	//gets two indexes in D
	cur_i, cur_j := rapidNeighborJoining(u, D, S, dead_records)

	//make sure p_i is the smallest index.
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
		fmt.Println(v_iu)
		fmt.Println(v_iu)

		distance_to_x, _ := strconv.ParseFloat(v_iu, 64)
		distance_to_y, _ := strconv.ParseFloat(v_ju, 64)

		newNode := integrateNewNode(array[cur_i], array[cur_j], distance_to_x, distance_to_y)
		array[cur_i] = newNode
		treeBanana = append(treeBanana, newNode)

		array = append(array[:cur_j], array[cur_j+1:]...)

		//creating newick form
		labels[cur_i] = "(" + labels[cur_i] + ":" + v_iu + "," + labels[cur_j] + ":" + v_ju + ")"
		labels = append(labels[:cur_j], labels[cur_j+1:]...)

		for i, v := range labels {
			fmt.Println(i, v, "This is god")
		}
	}

	D_new, S_new, dead_records_new := createNewDistanceMatrix(S, dead_records, D, cur_i, cur_j)

	fmt.Println("what is going")
	for i := 0; i < len(labels); i++ {
		fmt.Println(labels[i])
	}
	fmt.Println("okay nice")

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

	//assign dead records -> -1
	for k, v := range dead_records {
		if v == p_i {
			fmt.Println("also once per it")
			dead_records[k] = -1
		}
		if v == p_j {
			fmt.Println("once per it")
			dead_records[k] = -1
		}
		if v > p_j {
			dead_records[k] = v - 1
		}

	}
	dead_records[len(dead_records)] = p_i
	fmt.Println(dead_records)
	fmt.Println(D_new)

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
		var result int

		tuple.value = D[p_i][j]
		for k, v := range dead_records {
			if v == j {
				result = k
			}
		}
		tuple.index_j = result
		S_new[p_i][j] = tuple

	}
	//cut excess data away
	S_new[p_i] = S_new[p_i][:len(D_new)]

	//sort merged row
	sort.Slice(S_new[p_i], func(a, b int) bool {
		return (S_new[p_i][a].value < S_new[p_i][b].value)
	})

	fmt.Println("asdhf")
	for i := 0; i < len(S_new); i++ {
		fmt.Println(S_new[i])
	}

	//delete row in S
	S_new = append(S[:p_j], S[p_j+1:]...)

	return D_new, S_new, dead_records
}

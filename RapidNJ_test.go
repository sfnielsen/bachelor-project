package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func standardSetup(D [][]float64, labels []string) ([][]Tuple, map[int]int, Tree, Tree) {
	S := initSmatrix(D)
	deadRecords := initDeadRecords(D)
	var treeBanana Tree
	var array Tree
	array = generateTreeForRapidNJ(labels)

	for _, node := range array {
		treeBanana = append(treeBanana, node)
	}
	return S, deadRecords, array, treeBanana
}
func Test4Taxa(t *testing.T) {
	labels := []string{
		"A", "B", "C", "D",
	}
	D := [][]float64{
		{0, 17, 21, 27},
		{17, 0, 12, 18},
		{21, 12, 0, 14},
		{27, 18, 14, 0},
	}
	S, deadRecords, array, treeBanana := standardSetup(D, labels)

	newick_result, _ := neighborJoin(D, S, labels, deadRecords, array, treeBanana)
	print()
	if newick_result != "((B:4.000000,A:13.000000):2.000000,(C:4.000000,D:10.000000):2.000000);" {
		t.Errorf(newick_result)
	}

}

func Test4Taxa_made_up_numbers(t *testing.T) {
	labels := []string{
		"A", "B", "C", "D",
	}
	D := [][]float64{
		{0, 11, 18, 24},
		{11, 0, 12, 18},
		{18, 12, 0, 14},
		{24, 18, 14, 0},
	}
	S, deadRecords, array, treeBanana := standardSetup(D, labels)

	newick_result, _ := neighborJoin(D, S, labels, deadRecords, array, treeBanana)

	if newick_result != "((B:2.500000,A:8.500000):2.750000,(C:4.000000,D:10.000000):2.750000);" {
		t.Errorf("%s hehehe", newick_result)
	}

}

func Test8Taxa_madeUpNumbers_shouldBeChangedLater(t *testing.T) {

	labels := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	D := [][]float64{
		{0, 5, 68, 57, 127, 27, 28, 33},
		{5, 0, 58, 47, 117, 8, 52, 57},
		{68, 58, 0, 35, 69, 35, 87, 92},
		{57, 47, 35, 0, 94, 44, 79, 84},
		{127, 117, 69, 94, 0, 144, 149, 154},
		{27, 8, 35, 44, 144, 0, 27, 54},
		{28, 52, 87, 79, 149, 27, 0, 13},
		{33, 57, 92, 84, 154, 54, 13, 0},
	}

	S, deadRecords, array, treeBanana := standardSetup(D, labels)

	newick_result, _ := neighborJoin(D, S, labels, deadRecords, array, treeBanana)
	if newick_result != "(B:5.765625,((((G:1.250000,H:11.750000):23.208333,A:0.791667):4.718750,F:29.781250):1.281250,((C:0.333333,E:68.666667):18.200000,D:11.800000):2.968750):5.765625);" {
		t.Errorf("hehehe")
	}
}

func Test_max_taxa_of_generated_tree(t *testing.T) {
	prev_time := int64(0)
	quadratic := .0
	for i := 0; i < 5; i++ {

		taxa_amount := int(math.Pow(2, float64(i))) // power of 2
		time_start := time.Now().UnixMilli()
		generateTree(taxa_amount, 1)
		time_end := time.Now().UnixMilli()
		time := time_end - time_start

		if time != 0 && prev_time != 0 {
			quadratic = float64(time) / float64(prev_time*4) //suspected taxa ^ 2 running time
		}
		prev_time = time
		fmt.Printf("time for generating %d taxa is %d milliseconds  \nQuadratic time: %f.\n", taxa_amount, time, quadratic)
	}
}

func Test_Generated_Tree(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	taxa_amount := 51 + rand.Intn(51) //between 50 and 100
	tree, _, array := generateTree(taxa_amount, 10)

	//check if transposed distance matrix equals the distance matrix
	for i := range array {
		for j := range array {
			if i == j && array[i][j] != 0 {
				t.Errorf("diagonal not 0")
			}
			if array[i][j] != array[j][i] {
				t.Errorf("transpose not same as original")
			}
		}
	}

	//we are assuming that the tree indexes corresponds to the matrix indexes here
	//check if we can go through the tree and get same distance as written in the distance matrix
	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())
		//tree consists of 2n-2 nodes where n are leaves. We can only look at leaves. Note that start and to can be the same
		idx_start := rand.Intn(len(tree) / 2)
		idx_to := rand.Intn(len(tree) / 2)

		node_from := tree[idx_start]
		node_to_name := tree[idx_to].Name

		distance, _ := dfs_tree(node_from, node_to_name, make(map[*Node]bool))

		if distance != array[idx_start][idx_to] {
			t.Errorf("Distance should be the same. ")
		}
	}

	//test whether the nodes match the expected 2n+2, where n is taxa
	if count_nodes(tree[0], make(map[*Node]bool)) != (2*taxa_amount - 2) {
		t.Errorf("amount of nodes in tree does not fit 2n-2 as expected")
	}
	//tree should also contain 2n+2
	if len(tree) != (2*taxa_amount - 2) {
		t.Errorf("tree should contain excactly all taxa")
	}

}

func TestMakeTree(t *testing.T) {
	a, b, c := generateTree(5, 3)

	if a == nil || b == nil || c == nil {
		t.Errorf("poops")
	}
}

func TestRapidNJWithRandomDistanceMatrix(t *testing.T) {
	_, labels, distanceMatrix := generateTree(7, 20)
	original_labels := make([]string, len(labels))
	copy(original_labels, labels)

	original_dist_mat := make([][]float64, len(distanceMatrix))
	for i := range distanceMatrix {
		original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
		copy(original_dist_mat[i], distanceMatrix[i])
	}

	S, dead_record, array, treeBanana := standardSetup(distanceMatrix, labels)

	_, resulting_tree := neighborJoin(distanceMatrix, S, labels, dead_record, array, treeBanana)

	fmt.Println(len(resulting_tree))

	emptyMatrix := make([][]float64, len(original_labels))
	for i := range distanceMatrix {
		emptyMatrix[i] = make([]float64, len(original_labels))
	}

	fmt.Println("hello hello")
	for i := 0; i < len(original_dist_mat); i++ {
		fmt.Println(original_dist_mat[i])
	}

	fmt.Println("???")
	resulting_distance_matrix := createDistanceMatrix(emptyMatrix, resulting_tree, original_labels)
	for i := 0; i < len(resulting_distance_matrix); i++ {
		fmt.Println(resulting_distance_matrix[i])
	}
	are_they_the_same := compareDistanceMatrixes(distanceMatrix, resulting_distance_matrix)

	if !are_they_the_same {
		t.Errorf("failure :(")
	}

}

//#############################################
//helper functions we use in the test framework
//#############################################

//this is not used could perhaps be deleted

func compareDistanceMatrixes(matrix1 [][]float64, matrix2 [][]float64) bool {
	for i, row := range matrix1 {
		for j := range row {
			if matrix1[i][j] != matrix2[i][j] {
				return false
			}
		}
	}
	return true
}

func transposeMatrix(matrix [][]float64) [][]float64 {
	size := len(matrix)
	transposed := make([][]float64, size)
	for i := range transposed {
		transposed[i] = make([]float64, size)
	}

	for i, rows := range matrix {
		for j := range rows {
			transposed[j][i] = matrix[i][j]
		}

	}
	return transposed
}

//depth first searching on a tree of nodes starting at current_node. Note that -1 means that destionation was not found
func dfs_tree(current_node *Node, destination_name string, marked map[*Node]bool) (float64, *Node) {
	marked[current_node] = true
	distance := .0

	if current_node.Name == destination_name {
		return distance, current_node
	}
	for _, edge := range current_node.Edge_array {
		if _, ok := marked[edge.Node]; ok {
			continue
		}
		//check if we are looking at a leaf
		if len(edge.Node.Edge_array) == 1 {
			//check if leaf is the desired destionation
			if edge.Node.Name == destination_name {
				distance += edge.Distance
				return distance, current_node
			}

		} else {
			value, node := dfs_tree(edge.Node, destination_name, marked)
			if node != nil {
				distance = value + edge.Distance
				return distance, node
			}
		}
	}
	return -1, nil
}

//should take a node and traverse the tree the node is connected to. Returning the total amount of nodes in the tree. This should be 2*taxa-2
func count_nodes(current_node *Node, marked map[*Node]bool) (total_nodes int) {
	marked[current_node] = true
	sum := 1
	for _, edge := range current_node.Edge_array {
		if _, ok := marked[edge.Node]; ok {
			continue
		}

		//if edgearray we are going to only has one edge it must be a dead end
		if len(edge.Node.Edge_array) == 1 {
			sum++
		} else {
			sum += count_nodes(edge.Node, marked)
		}
	}
	return sum
}

func compare_trees(tree1 Tree, tree2 Tree) bool {
	var tree1_node *Node
	for _, v := range tree1 {
		//identify some actual taxa
		if len(v.Edge_array) == 1 {
			tree1_node = v
			break
		}
	}

	_, tree2_node := dfs_tree(tree2[0], tree1_node.Name, make(map[*Node]bool))

	//printing so that compiler does not complain xD
	fmt.Println(tree2_node)
	//############################
	//now the two identical taxa/tip nodes are found and we can start the proces of traversing the tree from this node and comparing.
	//############################
	return true
}

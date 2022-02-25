package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func standardSetup(D [][]float64) ([][]Tuple, map[int]int) {
	S := initSmatrix(D)
	deadRecords := initDeadRecords(D)
	return S, deadRecords
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
	S, deadRecords := standardSetup(D)

	newick_result := neighborJoin(D, S, labels, deadRecords)
	if newick_result != "((B:4.000000,A:13.000000):2.000000,(C:4.000000,D:10.000000):2.000000);" {
		t.Errorf("hehehe")
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
	S, deadRecords := standardSetup(D)

	newick_result := neighborJoin(D, S, labels, deadRecords)
	if newick_result != "((B:2.500000,A:8.500000):2.750000,(C:4.000000,D:10.000000):2.750000);" {
		t.Errorf("hehehe")
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

	S, deadRecords := standardSetup(D)

	newick_result := neighborJoin(D, S, labels, deadRecords)
	if newick_result != "(B:5.765625,((((G:1.250000,H:11.750000):23.208333,A:0.791667):4.718750,F:29.781250):1.281250,((C:0.333333,E:68.666667):18.200000,D:11.800000):2.968750):5.765625);" {
		t.Errorf("hehehe")
	}
}

func Test_Generated_Tree(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tree, _, array := generateTree(51+rand.Intn(50), 10)
	//check if transposed distance matrix equals the distance matrix
	for i := range array {
		for j := range array {
			if i == j && array[i][j] != 0 {
				fmt.Println(i, j)
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

		distance := dfs_tree(node_from, node_to_name, 0, make(map[*Node]bool))

		if distance != array[idx_start][idx_to] {
			t.Errorf("Distance should be the same. ")
		}
	}

}

func TestMakeTree(t *testing.T) {
	a, b, c := generateTree(5, 3)

	if a == nil || b == nil || c == nil {
		t.Errorf("poops")
	}
}

//#############################################
//helper functions we use in the test framework
//#############################################

//this is not used could perhaps be deleted
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
func dfs_tree(current_node *Node, destination_name string, sum float64, marked map[*Node]bool) (distance float64) {
	marked[current_node] = true

	if current_node.Name == destination_name {
		return
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
				return
			}

		} else {
			value := dfs_tree(edge.Node, destination_name, sum, marked)
			if value != -1 {
				return value + edge.Distance
			}
		}
	}
	return -1
}

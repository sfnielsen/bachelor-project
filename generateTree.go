package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Edge struct {
	node     *Node
	distance int
}

type Node struct {
	name string
	//if len(edge_array) == 1, the Node should be a 'leaf'
	edge_array []*Edge
}

type Tree []Node

func remove(slice []Node, s int) []Node {
	return append(slice[:s], slice[s+1:]...)
}

func generateTree(size int) (Tree, []string, [][]float64) {
	array := generateArray(size)
	tree := make(Tree, 0)

	//initialize distance matrix
	distanceMatrix := make([][]float64, size)
	for i := range distanceMatrix {
		distanceMatrix[i] = make([]float64, size)
	}
	//initialize labels
	labels := make([]string, size)

	//append all staring nodes to tree and create labels
	for _, value := range array {
		labels = append(labels, value.name)
		tree = append(tree, value)
	}
	distanceMatrix = createDistanceMatrix(distanceMatrix, labels, tree)

	fmt.Println("GOODDAY")
	for len(array) > 1 {

		fmt.Println("yo")
		for i := 0; i < len(array); i++ {
			fmt.Println(array[i].name)
		}
		fmt.Println("saka")

		rand.Seed(time.Now().UnixNano())

		random_x := rand.Intn(len(array))
		random_y := random_x

		//while loop that ensures we find two unique random integers
		for random_x == random_y {
			random_y = rand.Intn(len(array))
		}
		fmt.Println("swamp")
		fmt.Println("length of array:", len(array))
		fmt.Println("drawn")
		fmt.Println(random_x)
		fmt.Println(random_y)
		fmt.Println("oky")

		element_x := array[random_x]
		element_y := array[random_y]

		if random_x < random_y {
			array = remove(array, random_y)
			array = remove(array, random_x)
		} else {
			array = remove(array, random_x)
			array = remove(array, random_y)
		}

		//initialize new node and set its name as appended string combination
		new_node := new(Node)
		new_node.name = "(" + element_x.name + "," + element_y.name + ")"

		//make pointers to joined nodes
		new_edge_a := new(Edge)
		new_edge_a.distance = rand.Intn(20)
		new_edge_a.node = &element_x

		new_edge_b := new(Edge)
		new_edge_b.distance = rand.Intn(20)
		new_edge_b.node = &element_y

		//append to edges to new node's array
		new_node.edge_array = append(new_node.edge_array, new_edge_a)
		new_node.edge_array = append(new_node.edge_array, new_edge_b)

		//make pointers to new node
		edge_to_new_node_from_a := new(Edge)
		edge_to_new_node_from_a.distance = new_edge_a.distance
		edge_to_new_node_from_a.node = new_node

		edge_to_new_node_from_b := new(Edge)
		edge_to_new_node_from_b.distance = new_edge_a.distance
		edge_to_new_node_from_b.node = new_node

		//append edge to new node to joined neighbours' edge-arrays
		element_x.edge_array = append(element_x.edge_array, edge_to_new_node_from_a)
		element_y.edge_array = append(element_y.edge_array, edge_to_new_node_from_b)

		array = append(array, *new_node)

		tree = append(tree, *new_node)

	}
	fmt.Println(array[0].name)

	return tree, labels, distanceMatrix
}

func generateArray(numberOfLeafs int) []Node {

	returnArray := make([]Node, numberOfLeafs)
	fmt.Println("bevore lop")
	for i := 0; i < numberOfLeafs; i++ {
		fmt.Println("kek")
		node := new(Node)
		node.name = strconv.Itoa(i)
		fmt.Println(node.name)
		fmt.Println(returnArray)
		returnArray[i] = *node
		fmt.Println(returnArray)
	}
	return returnArray
}

//implement the sort.Interface interface for the Tree datatype
func (a Tree) Len() int           { return len(a) }
func (a Tree) Less(i, j int) bool { return a[i].name < a[j].name }
func (a Tree) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func createDistanceMatrix(distanceMatrix [][]float64, labels []string, tree Tree) [][]float64 {
	for _, node := range tree {
		if node.edge_array == 1 {

			//node.name

		}
	}
	return distanceMatrix
}

func sortTree() {

}

func compare_trees(tree1 Tree, tree2 Tree) bool {
	return true
}

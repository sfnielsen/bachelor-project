package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Edge struct {
	Node     *Node
	Distance int
}

type Node struct {
	Name string
	//if len(edge_array) == 1, the Node should be a 'leaf'
	Edge_array []*Edge
}

type Tree []*Node

func remove(slice []*Node, s int) []*Node {
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

	//append all staring nodes to tree and create labels
	labels := make([]string, 0)

	for _, value := range array {
		labels = append(labels, value.Name)
		tree = append(tree, value)
	}

	for len(array) > 1 {

		fmt.Println("yo")
		for i := 0; i < len(array); i++ {
			fmt.Println(array[i].Name)
		}

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
		new_node.Name = "(" + element_x.Name + "," + element_y.Name + ")"

		//make pointers to joined nodes
		new_edge_a := new(Edge)
		new_edge_a.Distance = rand.Intn(20)
		new_edge_a.Node = element_x

		new_edge_b := new(Edge)
		new_edge_b.Distance = rand.Intn(20)
		new_edge_b.Node = element_y

		//append to edges to new node's array
		new_node.Edge_array = append(new_node.Edge_array, new_edge_a)
		new_node.Edge_array = append(new_node.Edge_array, new_edge_b)

		//make pointers to new node
		edge_to_new_node_from_a := new(Edge)
		edge_to_new_node_from_a.Distance = new_edge_a.Distance
		edge_to_new_node_from_a.Node = new_node

		edge_to_new_node_from_b := new(Edge)
		edge_to_new_node_from_b.Distance = new_edge_a.Distance
		edge_to_new_node_from_b.Node = new_node
		fmt.Println("swpm")

		fmt.Println(len(element_x.Edge_array))

		//append edge to new node to joined neighbours' edge-arrays
		element_x.Edge_array = append(element_x.Edge_array, edge_to_new_node_from_a)
		element_y.Edge_array = append(element_y.Edge_array, edge_to_new_node_from_b)

		fmt.Println(len(element_x.Edge_array))
		array = append(array, new_node)

		tree = append(tree, new_node)

	}
	fmt.Println(array[0].Name)

	distanceMatrix = createDistanceMatrix(distanceMatrix, tree, labels)
	return tree, labels, distanceMatrix
}

func generateArray(numberOfLeafs int) []*Node {

	returnArray := make([]*Node, numberOfLeafs)
	fmt.Println("bevore lop")
	for i := 0; i < numberOfLeafs; i++ {
		fmt.Println("kek")
		node := new(Node)
		node.Name = strconv.Itoa(i)
		fmt.Println(node.Name)
		fmt.Println(returnArray)
		returnArray[i] = node
		fmt.Println(returnArray)
	}
	return returnArray
}

//implement the sort.Interface interface for the Tree datatype
func (a Tree) Len() int           { return len(a) }
func (a Tree) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a Tree) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func traverseTree(distanceRow []float64, node Node, sum float64, seen map[string]bool, labelMap map[string]int) []float64 {
	if _, ok := seen[node.Name]; ok {
		return distanceRow
	}
	//set THIS node to seen
	seen[node.Name] = true
	if len(node.Edge_array) == 1 {
		distanceRow[labelMap[node.Name]] = sum
		return distanceRow
	}

	for _, edge := range node.Edge_array {
		sum1 := sum
		sum1 += float64(edge.Distance)
		distanceRow = traverseTree(distanceRow, *edge.Node, sum1, seen, labelMap)
	}
	return distanceRow
}

func createDistanceMatrix(distanceMatrix [][]float64, tree Tree, labels []string) [][]float64 {
	print("creating the distance")
	labelMap := make(map[string]int)
	for i, v := range labels {
		labelMap[v] = i
	}
	fmt.Println(len(labelMap))

	for _, node := range tree {
		fmt.Println("investigating the future")
		fmt.Println(len(node.Edge_array))
		if len(node.Edge_array) == 1 {
			seen := make(map[string]bool)
			fmt.Println("the length is the one")

			distanceRow := make([]float64, len(labels))
			fmt.Println(len(distanceRow))
			fmt.Println(len(distanceMatrix))

			//this assumes that index 0 in array holds the lexicographicly first node. Perhaps sorting should be implemented to ensure this property
			distanceMatrix[labelMap[node.Name]] = traverseTree(distanceRow, *node.Edge_array[0].Node,
				float64(node.Edge_array[0].Distance), seen, labelMap)

			fmt.Println(distanceMatrix[labelMap[node.Name]])
		}
	}
	return distanceMatrix
}

func sortTree() {

}

func compare_trees(tree1 Tree, tree2 Tree) bool {
	return true
}

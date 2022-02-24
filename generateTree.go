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

func generateTree(size int, max_length_random int) (Tree, []string, [][]float64) {
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
		new_edge_a.Distance = rand.Intn(max_length_random) + 1
		new_edge_a.Node = element_x

		new_edge_b := new(Edge)
		new_edge_b.Distance = rand.Intn(max_length_random) + 1
		new_edge_b.Node = element_y

		//append to edges to new node's array
		new_node.Edge_array = append(new_node.Edge_array, new_edge_a)
		new_node.Edge_array = append(new_node.Edge_array, new_edge_b)

		//make pointers to new node
		edge_to_new_node_from_a := new(Edge)
		edge_to_new_node_from_a.Distance = new_edge_a.Distance
		edge_to_new_node_from_a.Node = new_node

		edge_to_new_node_from_b := new(Edge)
		edge_to_new_node_from_b.Distance = new_edge_b.Distance
		edge_to_new_node_from_b.Node = new_node

		//append edge to new node to joined neighbours' edge-arrays
		element_x.Edge_array = append(element_x.Edge_array, edge_to_new_node_from_a)
		element_y.Edge_array = append(element_y.Edge_array, edge_to_new_node_from_b)

		array = append(array, new_node)

		tree = append(tree, new_node)

		for i := 0; i < len(array); i++ {
			fmt.Println(array[i].Name)
		}

		//joining the last 2 nodes
		if len(array) == 2 {
			fmt.Println("yoda")
			//index 1 must be the one we just joined. We want to merge index 0 into this one aswell.

			array[1].Name = "(" + array[0].Name + "," + array[1].Name[1:]

			dist := rand.Intn(max_length_random) + 1

			fmt.Println("last dist", dist)
			new_edge_0 := new(Edge)
			new_edge_0.Distance = dist
			new_edge_0.Node = array[1]

			new_edge_1 := new(Edge)
			new_edge_1.Distance = dist
			new_edge_1.Node = array[0]

			array[0].Edge_array = append(array[0].Edge_array, new_edge_0)
			array[1].Edge_array = append(array[1].Edge_array, new_edge_1)

			array = remove(array, 0)
		}
	}

	fmt.Println(array[0].Name)

	distanceMatrix = createDistanceMatrix(distanceMatrix, tree, labels)
	return tree, labels, distanceMatrix
}

func generateArray(numberOfLeafs int) []*Node {

	returnArray := make([]*Node, numberOfLeafs)

	for i := 0; i < numberOfLeafs; i++ {

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
		new_sum := sum
		new_sum += float64(edge.Distance)
		distanceRow = traverseTree(distanceRow, *edge.Node, new_sum, seen, labelMap)
	}
	return distanceRow
}

func createDistanceMatrix(distanceMatrix [][]float64, tree Tree, labels []string) [][]float64 {
	print("creating the distance\n")
	labelMap := make(map[string]int)
	for i, v := range labels {
		labelMap[v] = i
	}

	for _, node := range tree {

		if len(node.Edge_array) == 1 {
			//initialize seen map (set) and adding the current node
			seen := make(map[string]bool)
			seen[node.Name] = true

			distanceRow := make([]float64, len(labels))

			//this assumes that index 0 in array holds the lexicographicly first node. Perhaps sorting should be implemented to ensure this property
			//we start from the only node that our current label connects to.
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

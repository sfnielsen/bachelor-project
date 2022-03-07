package main

import (
	"math/rand"
	"strconv"
	"time"
)

//Edge is a structure stored in a Node. It keeps tracks of the distance and a destination Node (different from the Node it is stored in)
type Edge struct {
	Node     *Node
	Distance float64
}

//Node is what makes up the phylegenetic trees. Every taxa is represented by a Node as well as every intersection between edges.
type Node struct {
	Name string
	//if len(edge_array) == 1, we consider the node a 'leaf/tip'
	Edge_array []*Edge
}

//the idea of the Tree structure is that it hold all nodes, both the labelled and interconnecting nodes
//for a given phylegenetic tree.
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

	//maybe it is not necessary to initialize new array makybe try without
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

		element_x := array[random_x]
		element_y := array[random_y]

		if random_x < random_y {
			array = remove(array, random_y)
			array = remove(array, random_x)
		} else {
			array = remove(array, random_x)
			array = remove(array, random_y)
		}

		distance_to_x := float64(rand.Intn(max_length_random) + 1)
		distance_to_y := float64(rand.Intn(max_length_random) + 1)

		var new_node *Node = integrateNewNode(element_x, element_y, distance_to_x, distance_to_y)

		// we need both - array holds only the last one in the end, tree holds every node
		array = append(array, new_node)

		tree = append(tree, new_node)

		//joining the last 2 nodes
		if len(array) == 2 {
			//index 1 must be the one we just joined. We want to merge index 0 into this one aswell.

			array[1].Name = "(" + array[0].Name + "," + array[1].Name[1:]

			dist := float64(rand.Intn(max_length_random)) + 1.0

			new_edge_0 := new(Edge)
			new_edge_0.Distance = float64(dist)
			new_edge_0.Node = array[1]

			new_edge_1 := new(Edge)
			new_edge_1.Distance = float64(dist)
			new_edge_1.Node = array[0]

			array[0].Edge_array = append(array[0].Edge_array, new_edge_0)
			array[1].Edge_array = append(array[1].Edge_array, new_edge_1)

			array = remove(array, 0)
		}
	}

	distanceMatrix = createDistanceMatrix(distanceMatrix, tree, labels)
	return tree, labels, distanceMatrix
}

func integrateNewNode(element_x *Node, element_y *Node, distance_to_x float64, distance_to_y float64) *Node {
	//initialize new node and set its name as appended string combination
	new_node := new(Node)
	new_node.Name = "(" + element_x.Name + "," + element_y.Name + ")"

	//make pointers to joined nodes
	new_edge_a := new(Edge)
	new_edge_a.Distance = distance_to_x
	new_edge_a.Node = element_x

	new_edge_b := new(Edge)
	new_edge_b.Distance = distance_to_y
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

	return new_node
}

func generateArray(numberOfLeafs int) []*Node {

	returnArray := make([]*Node, numberOfLeafs)

	for i := 0; i < numberOfLeafs; i++ {

		node := new(Node)
		node.Name = strconv.Itoa(i)
		returnArray[i] = node
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

	labelMap := make(map[string]int)

	//we need this to make distance matrix rows, so every row is in same order of labels.
	for i, v := range labels {
		labelMap[v] = i
	}

	for _, node := range tree {

		if len(node.Edge_array) == 1 {
			//initialize seen map (set) and adding the current node
			seen := make(map[string]bool)
			seen[node.Name] = true

			distanceRow := make([]float64, len(labels))

			//this assumes that index 0 in array holds the lexicographicly frst node. Perhaps sorting should be implemented to ensure this property
			//we start from the only node that our current label connects to.i
			distanceMatrix[labelMap[node.Name]] = traverseTree(distanceRow, *node.Edge_array[0].Node,
				float64(node.Edge_array[0].Distance), seen, labelMap)
		}
	}
	return distanceMatrix
}

func sortTree() {

}

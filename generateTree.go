package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Edge struct {
	node     *Node
	distance int
}

type Node struct {
	name       string
	edge_array []*Edge
}

func remove(slice []Node, s int) []Node {
	return append(slice[:s], slice[s+1:]...)
}

func generateTree() {
	array := generateArray(100)
	fmt.Println("GOODDAY")
	for len(array) > 3 {
		random_x := rand.Intn(len(array) - 1)
		random_y := random_x

		//while loop that ensures we find two unique random integers
		for random_x == random_y {
			random_y = rand.Intn(len(array) - 1)
		}

		element_x := array[random_x]
		element_y := array[random_y]

		array = remove(array, random_x)
		array = remove(array, random_y)

		//initialize new node and set its name as appended string combination
		new_node := new(Node)
		new_node.name = element_x.name + ":" + element_y.name

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
		fmt.Println("ugabuga")
		for i := 0; i > len(array)-1; i++ {
			print(array[i].name)
		}
	}
}

func generateArray(numberOfLeafs int) []Node {

	returnArray := make([]Node, numberOfLeafs)
	for i := 0; i > numberOfLeafs; i++ {
		leaf := new(Node)
		leaf.name = strconv.Itoa(i)
		returnArray = append(returnArray, *leaf)
	}
	return returnArray
}

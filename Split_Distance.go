package main

import (
	"reflect"
)

//The first node is from the first tree, the second node is the connection to the second tree
//first nodes' label is then searched for in second tree and this label is used as the root for the search
//Returns 0 if trees have the same splits.
//This function assumes two trees with the same labels are recieved as input.
func Split_Distance(node1 *Node, node2 *Node) int {

	//assume node from tree 1 is root
	node1_root := node1
	root := node1_root.Name
	//find declared root in second tree
	node2_root := node2
	if root != node2.Name {
		_, node2_root = dfs_tree(node2, root, make(map[*Node]bool))
	}
	tree1_splits := find_splits(node1_root, make([]*split, 0), root)
	tree2_splits := find_splits(node2_root, make([]*split, 0), root)

	difference := compare_trees(tree1_splits, tree2_splits)

	return difference
}

func compare_trees(splits_tree1 []*split, splits_tree2 []*split) int {

	//it is not easy to dynamically delete from split1 when we iterate through it. We use a variable to keep count of remaining instead.
	splits1_remaining := len(splits_tree1)

	for _, split1 := range splits_tree1 {

		for i, split2 := range splits_tree2 {
			//check if the head and tail are the same

			if reflect.DeepEqual(split1.head, split2.head) && reflect.DeepEqual(split1.tail, split2.tail) {
				splits_tree2 = append(splits_tree2[:i], splits_tree2[i+1:]...)
				splits1_remaining--
				break
			}

		}

	}
	unmatched_splits := splits1_remaining + len(splits_tree2)
	return unmatched_splits

}

//An edge in the some tree has 2 'Edge' since the edges depends on
//which node we go through
type split struct {
	loc1 *Edge
	loc2 *Edge
	head map[string]bool
	tail map[string]bool
}

//takes an empty map and recursively traverse the given Node
//in order to find all possible splits in the tree
func find_splits(current_node *Node, current_splits []*split, root string) []*split {
	for _, edge := range current_node.Edge_array {
		//check if some edge is already traversed
		traversed := false
		for _, split := range current_splits {
			if edge == split.loc1 || edge == split.loc2 {
				traversed = true
			}
		}
		if traversed {
			continue
		}

		//if not already traversed we create a new split
		to_node := edge.Node
		for _, to_edge := range to_node.Edge_array {
			if current_node.Name == to_edge.Node.Name {
				new_split := new(split)
				new_split.loc1 = edge
				new_split.loc2 = to_edge
				new_split.head, new_split.tail = split_tree(edge, to_edge, root)
				current_splits = append(current_splits, new_split)
				current_splits = find_splits(to_node, current_splits, root)
				break
			}
		}
	}
	return current_splits
}

//returns 'head' (where root is located) as first parameter and the other subtree 'tail'
//as the second parameter
func split_tree(loc1 *Edge, loc2 *Edge, root string) (head map[string]bool, tail map[string]bool) {

	//add the node from the other split to initial seen map to exclude the subtree
	loc1_marked := make(map[*Node]bool)
	loc1_marked[loc2.Node] = true
	loc2_marked := make(map[*Node]bool)
	loc2_marked[loc1.Node] = true

	_, head = find_subtree_nodes(loc1.Node, loc1_marked, make(map[string]bool))
	_, tail = find_subtree_nodes(loc2.Node, loc2_marked, make(map[string]bool))

	//make sure head is the subtree where the root is located
	if _, ok := tail[root]; ok {
		a := tail
		tail = head
		head = a
	}
	return head, tail
}

//Border is the edge that cuts off rest of the tree from the current subtree
//returns the marked map (used for recursive calls) and the labels found in subtree
func find_subtree_nodes(current_node *Node, marked map[*Node]bool, subtree_Nodes map[string]bool) (map[*Node]bool, map[string]bool) {
	if len(current_node.Edge_array) == 1 {
		subtree_Nodes[current_node.Name] = true
	}
	marked[current_node] = true
	for _, edge := range current_node.Edge_array {
		if _, ok := marked[edge.Node]; ok {
			continue
		}
		marked, subtree_Nodes = find_subtree_nodes(edge.Node, marked, subtree_Nodes)

	}
	return marked, subtree_Nodes
}

//depth first searching on a tree of nodes starting at current_node. Note that -1 means that destionation was not found.
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

package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

//####################################################################################
//####################################################################################
//####################################################################################
//                               testing other files
//####################################################################################
//####################################################################################
//####################################################################################

func TestMakeTree(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	a, b, c := GenerateTree(5, 3, Uniform_distribution, seed)

	if a == nil || b == nil || c == nil {
		t.Errorf("not good")
	}
}

func Test_max_taxa_of_generated_tree(t *testing.T) {
	prev_time := int64(0)
	quadratic := .0
	for i := 0; i < 5; i++ {
		seed := time.Now().UTC().UnixNano()
		taxa_amount := int(math.Pow(2, float64(i))) // power of 2
		time_start := time.Now().UnixMilli()
		GenerateTree(taxa_amount, 1, Normal_distribution, seed)
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
	seed := time.Now().UnixNano()
	taxa_amount := 51 + rand.Intn(51) //between 50 and 100
	tree, _, array := GenerateTree(taxa_amount, 5, Normal_distribution, seed)

	//check if transposed distance matrix equals the distance matrix
	for i := range array {
		for j := range array {
			if i == j && array[i][j] != 0 {
				t.Errorf("diagonal not 0")
			}
			//account for float pointer precision 2 decimals
			round_idx1, round_idx2 := math.Round(array[i][j]*100)/100, math.Round(array[j][i]*100)/100
			if round_idx1 != round_idx2 {
				t.Errorf("transpose not same as original")
			}
		}
	}

	//we are assuming that the tree indexes corresponds to the matrix indexes here
	//check if we can go through the tree and get same distance as written in the distance matrix
	wrong_numbers := 0
	for i := 0; i < 1000; i++ {
		rand.Seed(time.Now().UnixNano())
		//tree consists of 2n-2 nodes where n are leaves. We can only look at leaves. Note that start and to can be the same
		idx_start := rand.Intn(len(tree) / 2)
		idx_to := rand.Intn(len(tree) / 2)

		node_from := tree[idx_start]
		node_to_name := tree[idx_to].Name

		distance, _ := dfs_tree(node_from, node_to_name, make(map[*Node]bool))

		//account for float pointer precision 2 decimals
		round_dist, round_exp_dist := math.Round(distance*100/100), math.Round(array[idx_start][idx_to]*100/100)
		if round_dist != round_exp_dist {
			wrong_numbers++
			//some times the floats are off by a little much so we allow a couple of them to be off
			if wrong_numbers > 10 {
				fmt.Println(idx_start, idx_to)
				fmt.Println(round_dist, round_exp_dist)

				t.Errorf("Distance should be the same. ")

			}
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

func Test_Split_Distance(t *testing.T) {
	iterations := 50
	results := make(map[int]float64)
	seed := time.Now().UTC().UnixNano()
	for i := 0; i < iterations; i++ {

		//GENERATE 2 TREES
		tree1, _, _ := GenerateTree(5, 15, Normal_distribution, seed)
		tree2, _, _ := GenerateTree(5, 15, Normal_distribution, seed+1)

		//CHECK TREES
		results[Split_Distance(tree1[0], tree2[2])]++
	}

	for k, v := range results {
		results[k] = float64(v) / float64(iterations)
	}

	//WE EXPECT 1/15 OF RANDOM 5-tip TREES TO BE TOPOLOGICALLY IDENTICAL
	test := math.Abs(results[0] - float64(1)/float64(15))
	if test > 0.1 {
		t.Errorf("Expect 1/15 good trees, %f", test)
	}

	//CHECK THAT BIGGER TREES DO NOT RANDOMLY MAKE THE SAME TREE
	for i := 0; i < (iterations / 5); i++ {

		//GENERATE 3 TREES
		tree1, _, _ := GenerateTree(50, 15, Normal_distribution, seed)
		tree2, _, _ := GenerateTree(50, 15, Normal_distribution, seed+1)

		tree1_identical, _, _ := GenerateTree(50, 15, Normal_distribution, seed)

		//CHECK TREES
		result := Split_Distance(tree1[0], tree2[2])
		if result == 0 {
			t.Errorf("Unlikely scenario - Big trees not expected to randomly be identical")
		}
		result2 := Split_Distance(tree1[0], tree1_identical[2])
		if result2 != 0 {
			t.Error("These two trees should have distance 0 - eg. being the same")
		}
	}
}

func Test_Split_Distance_fails(t *testing.T) {

	taxa := 100
	seed := time.Now().UTC().UnixNano()
	_, labels1, distanceMatrix1 := GenerateTree(taxa/2, 15, Normal_distribution, seed)

	//trees have different amount of taxas such that all splits beocome different

	_, labels2, distanceMatrix2 := GenerateTree(taxa, 15, Normal_distribution, seed)

	//_, _, array, tree := standardSetup(distanceMatrix1, labels1)
	_, canon_tree := neighborJoin(distanceMatrix1, labels1)

	_, rapid_tree := rapidNeighbourJoin(distanceMatrix2, labels2, rapidNeighborJoining)

	test := Split_Distance(canon_tree[0], rapid_tree[0])
	fmt.Println(test)
	//since graph always has 1 less edge than node we expect errors to be n*2-3, and n-3 for the other one.
	if test != (taxa*3 - 3*2) {
		t.Errorf("split distance is 0! Trees should not be the same")
	}
}

//####################################################################################
//####################################################################################
//####################################################################################
//                               testing of RapidNJ
//####################################################################################
//####################################################################################
//####################################################################################

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

	newick_result, _ := rapidNeighbourJoin(D, labels, rapidNeighborJoining)

	//note that the newick always becomes a rooted tree whereas our implementation of the algorithm generates an unrooted tree.
	if newick_result != "(((A:13.000000,B:4.000000):4.000000,C:4.000000):5.000000,D:5.000000);" {
		t.Errorf(newick_result)
	}

}

func TestRapidNJ20TaxaRandomDistMatrix100Times(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	for i := 0; i < 100; i++ {
		seed++
		_, labels, distanceMatrix := GenerateTree(100, 5, Uniform_distribution, seed)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)

		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		_, resulting_tree := rapidNeighbourJoin(distanceMatrix, labels, rapidNeighborJoining)
		emptyMatrix := make([][]float64, len(labels))
		for i := range distanceMatrix {
			emptyMatrix[i] = make([]float64, len(labels))
		}
		resulting_distance_matrix := createDistanceMatrix(emptyMatrix, resulting_tree, original_labels)
		are_they_the_same := compareDistanceMatrixes(original_dist_mat, resulting_distance_matrix)

		if !are_they_the_same {
			t.Errorf(" failure :(")
		}
	}

}

func Test_Profiling_on_rapidNeighbourJoin(t *testing.T) {
	fmt.Println("...running profiling...")

	taxa := 2000

	var time_start, time_end, time_measured_rapid int64

	NewickFlag = true
	//seed := time.Now().UTC().UnixNano()
	seed := int64(6969)
	_, labels, distanceMatrix := GenerateTree(taxa, 15, Normal_distribution, seed)
	time_start = time.Now().UnixMilli()
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
	rapidNeighbourJoin(distanceMatrix, labels, rapidNeighborJoining)
	pprof.StopCPUProfile()

	time_end = time.Now().UnixMilli()
	time_measured_rapid = time_end - time_start
	fmt.Printf("### TIME ELAPSED: %d ms for run ###\n", time_measured_rapid)
	fmt.Printf("### TOTAL AMOUNT OF LOOKUPS: %d ###\n", total_lookups)
	fmt.Printf("### COLUMN UPDATES WAS FOUND ###\n")

	//fmt.Println(column_depth)
	//fmt.Println(total_updates)
	//fmt.Println(extra_cost)
	//fmt.Println(total)

	if 1 == 2 {
		t.Errorf("error")
	}
}

func TestRapidNJWithRandomDistanceMatrix(t *testing.T) {
	NewickFlag = true
	seed := time.Now().UTC().UnixNano()
	for i := 0; i < 1; i++ {
		seed++
		_, labels, distanceMatrix := GenerateTree(1500, 40, Normal_distribution, seed)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)

		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		fmt.Println("###DO NEIGHBOURJOIN")
		_, resulting_tree := rapidNeighbourJoin(distanceMatrix, labels, rapidNeighborJoining)

		emptyMatrix := make([][]float64, len(labels))
		fmt.Println("###CREATE DISTANCE MATRIX")
		for i := range distanceMatrix {
			emptyMatrix[i] = make([]float64, len(original_labels))
		}

		fmt.Println("###ORIGINAL MATRIX")

		//for i := 0; i < len(original_dist_mat); i++ {
		//	fmt.Println(original_dist_mat[i])
		//}

		fmt.Println("###NOW CREATE DIST")
		resulting_distance_matrix := createDistanceMatrix(emptyMatrix, resulting_tree, original_labels)
		//for i := 0; i < len(resulting_distance_matrix); i++ {
		//	fmt.Println(resulting_distance_matrix[i])
		//}
		fmt.Println("###COMPARE WITH ORIGINAL")

		are_they_the_same := compareDistanceMatrixes(original_dist_mat, resulting_distance_matrix)
		if !are_they_the_same {
			t.Errorf(" failure :(")
		}
	}

}

func TestCanonicalNJ20TaxaRandomDistMatrix100Times(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	for i := 0; i < 100; i++ {
		seed++
		_, labels, distanceMatrix := GenerateTree(20, 15, Normal_distribution, seed)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)

		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		//_, _, array, tree := standardSetup(distanceMatrix, labels)
		_, resulting_tree := neighborJoin(distanceMatrix, labels)
		emptyMatrix := make([][]float64, len(labels))
		for i := range distanceMatrix {
			emptyMatrix[i] = make([]float64, len(labels))
		}

		resulting_distance_matrix := createDistanceMatrix(emptyMatrix, resulting_tree, original_labels)
		are_they_the_same := compareDistanceMatrixes(original_dist_mat, resulting_distance_matrix)

		if !are_they_the_same {
			t.Errorf(" failure :(")
		}
	}
}

func Test_Canonical_rapid_generate_identical_matrixes_and_split_distance0(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	_, labels, distanceMatrix := GenerateTree(500, 15, Normal_distribution, seed)
	original_labels := make([]string, len(labels))
	copy(original_labels, labels)

	original_dist_mat := make([][]float64, len(distanceMatrix))
	for i := range distanceMatrix {
		original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
		copy(original_dist_mat[i], distanceMatrix[i])
	}

	//_, _, array, tree := standardSetup(distanceMatrix, labels)
	_, canon_tree := neighborJoin(distanceMatrix, labels)
	emptyMatrix1 := make([][]float64, len(labels))
	for i := range distanceMatrix {
		emptyMatrix1[i] = make([]float64, len(labels))
	}

	resulting_canonical_matrix := createDistanceMatrix(emptyMatrix1, canon_tree, original_labels)

	labels_cpy := make([]string, len(original_labels))
	copy(labels_cpy, original_labels)

	dist_mat_cpy := make([][]float64, len(original_dist_mat))
	for i := range original_dist_mat {
		dist_mat_cpy[i] = make([]float64, len(original_dist_mat[i]))
		copy(dist_mat_cpy[i], original_dist_mat[i])
	}

	_, rapid_tree := rapidNeighbourJoin(dist_mat_cpy, labels_cpy, rapidNeighborJoining)
	emptyMatrix2 := make([][]float64, len(labels_cpy))
	for i := range dist_mat_cpy {
		emptyMatrix2[i] = make([]float64, len(labels_cpy))
	}

	resulting_rapid_matrix := createDistanceMatrix(emptyMatrix2, rapid_tree, original_labels)

	cmp_canon_original := compareDistanceMatrixes(original_dist_mat, resulting_canonical_matrix)
	cmp_rapid_original := compareDistanceMatrixes(original_dist_mat, resulting_rapid_matrix)

	if !(cmp_canon_original) {
		t.Errorf("error: canon != origninal")
	}
	if !(cmp_rapid_original) {
		t.Errorf("error: rapid != original")
	}
	cmp_canon_rapid := compareDistanceMatrixes(resulting_canonical_matrix, resulting_rapid_matrix)
	//this case should not be possible of passed the two other comparisons
	if !(cmp_canon_rapid) {
		t.Errorf("error: canon != rapid")
	}

	errors := Split_Distance(canon_tree[0], rapid_tree[0])
	if errors != 0 {
		t.Errorf("error: rapid and canonical tree are not the same, %d errors", errors)
	}

}

func Test_Parse_phylip_distance_form_real_data_96_taxa(t *testing.T) {
	//parse phylip distance matrix format
	D1, labels1 := Parse_text()
	D2, labels2 := Parse_text()
	D_cpy, _ := Parse_text()
	//copy labels
	original_labels1 := make([]string, len(labels1))
	original_labels2 := make([]string, len(labels1))
	copy(original_labels1, labels1)
	copy(original_labels2, labels2)
	//neighbour joining (rapid and canonical)
	_, tree1 := rapidNeighbourJoin(D1, labels1, rapidNeighborJoining)
	_, tree2 := neighborJoin(D2, labels2)
	//split distance
	dist := Split_Distance(tree1[1], tree2[2])
	if dist != 0 {
		t.Errorf("not 0 split distance")
	}
	//distance matrix same
	d_new1 := make([][]float64, len(D1))
	d_new2 := make([][]float64, len(D1))
	for i := range D1 {
		d_new1[i] = make([]float64, len(D1))
		d_new2[i] = make([]float64, len(D1))
	}
	res_D1 := createDistanceMatrix(d_new1, tree1, original_labels1)
	res_D2 := createDistanceMatrix(d_new1, tree2, original_labels2)
	same_matrix1 := compareDistanceMatrixes(res_D1, res_D2)
	same_matrix2 := compareDistanceMatrixes(D_cpy, res_D1)
	if !(same_matrix1 && same_matrix2) {
		t.Errorf("Matrix are not the same")
	}
}

func TestUMAXheuristic(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	for i := 0; i < 100; i++ {
		seed++
		_, labels, distanceMatrix := GenerateTree(20, 15, Uniform_distribution, seed)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)

		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		_, resulting_tree := rapidNeighbourJoin(distanceMatrix, labels, rapidNeighborJoining)
		emptyMatrix := make([][]float64, len(labels))
		for i := range distanceMatrix {
			emptyMatrix[i] = make([]float64, len(labels))
		}
		resulting_distance_matrix := createDistanceMatrix(emptyMatrix, resulting_tree, original_labels)
		are_they_the_same := compareDistanceMatrixes(original_dist_mat, resulting_distance_matrix)

		if !are_they_the_same {
			t.Errorf(" failure :(")
		}
	}

}

//####################################################################################
//####################################################################################
//####################################################################################
//                 helper functions we use in the test framework
//####################################################################################
//####################################################################################
//####################################################################################

func compareDistanceMatrixes(matrix1 [][]float64, matrix2 [][]float64) bool {
	var errs int = 0
	for i, row := range matrix1 {
		for j := range row {
			//check if number deviates from eachother more than 15 %
			if math.Abs(matrix1[i][j]-matrix2[i][j]) > math.Max(matrix1[i][j], matrix2[i][j])*0.15 {
				errs++
				//fmt.Println("indexes:", i, j)
				//fmt.Println("first", matrix1[i][j])
				//fmt.Println("second", matrix2[i][j])

				//extreme case: if the small number is less than 60% of the big number we do not accept:
				if math.Abs(matrix1[i][j]-matrix2[i][j]) > math.Max(matrix1[i][j], matrix2[i][j])*0.6 {
					return false
				}
			}
		}
	}
	percentage_off := float64(errs) / math.Pow(float64(len(matrix1)), 2)
	return percentage_off < 0.10

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

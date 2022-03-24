package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	Shifting_Normal_Distribution string = "Sh_norm"
	Normal_distribution          string = "Norm"
	Uniform_distribution         string = "Uniform"
	Cluster_Normal_Distribution  string = "Cluster_norm"
	Spike_Normal_distribution    string = "Spike_norm"
)

func standardSetup(D [][]float64, labels []string) ([][]Tuple, map[int]int, Tree, Tree) {
	S := initSmatrix(D)
	deadRecords := initDeadRecords(D)
	var tree Tree
	var label_tree Tree = generateTreeForRapidNJ(labels)

	tree = append(tree, label_tree...)

	return S, deadRecords, label_tree, tree
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

	newick_result, _ := rapidJoin(D, S, labels, deadRecords, array, treeBanana)
	print()

	//note that the newick always becomes a rooted tree whereas our implementation of the algorithm generates an unrooted tree.
	if newick_result != "(((A:13.000000,B:4.000000):4.000000,C:4.000000):5.000000,D:5.000000);" {
		t.Errorf(newick_result)
	}

}

func Test_max_taxa_of_generated_tree(t *testing.T) {
	prev_time := int64(0)
	quadratic := .0
	for i := 0; i < 5; i++ {

		taxa_amount := int(math.Pow(2, float64(i))) // power of 2
		time_start := time.Now().UnixMilli()
		generateTree(taxa_amount, 1, Normal_distribution)
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
	tree, _, array := generateTree(taxa_amount, 5, Normal_distribution)

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

func TestMakeTree(t *testing.T) {
	a, b, c := generateTree(5, 3, Uniform_distribution)

	if a == nil || b == nil || c == nil {
		t.Errorf("not good")
	}
}

func TestRapidNJ20TaxaRandomDistMatrix100Times(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, labels, distanceMatrix := generateTree(20, 15, Uniform_distribution)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)

		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		S, dead_record, array, treeBanana := standardSetup(distanceMatrix, labels)
		_, resulting_tree := rapidJoin(distanceMatrix, S, labels, dead_record, array, treeBanana)
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
func TestRapidNJWithRandomDistanceMatrix(t *testing.T) {
	for i := 0; i < 1; i++ {
		_, labels, distanceMatrix := generateTree(100, 100, Spike_Normal_distribution)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)

		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		S, dead_record, array, tree := standardSetup(distanceMatrix, labels)

		fmt.Println("###DO NEIGHBOURJOIN")
		_, resulting_tree := rapidJoin(distanceMatrix, S, labels, dead_record, array, tree)

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

func TestRuntimeOfBigTaxas(t *testing.T) {
	//declaring some variables to hold times
	var time_start, time_end int64

	fmt.Println("###GENERATING DISTANCE MATRIX")
	time_start = time.Now().UnixMilli()
	_, labels, distanceMatrix := generateTree(1000, 1000, Normal_distribution)
	time_end = time.Now().UnixMilli()
	time_generateTree := time_end - time_start
	fmt.Printf("###Done in %d milliseconds\n", time_generateTree)

	S, dead_record, array, treeBanana := standardSetup(distanceMatrix, labels)

	fmt.Println("###BEGIN NEIGHBOR-JOINING")
	time_start = time.Now().UnixMilli()
	a, b := rapidJoin(distanceMatrix, S, labels, dead_record, array, treeBanana)
	time_end = time.Now().UnixMilli()
	time_neighborJoin := int(time_end - time_start)
	fmt.Printf("###Done in %d milliseconds\n", time_neighborJoin)

	if a == "" || b == nil {
		t.Errorf(" failure :(")
	}
}

func TestCanonicalNJ20TaxaRandomDistMatrix100Times(t *testing.T) {
	for i := 0; i < 100; i++ {
		_, labels, distanceMatrix := generateTree(20, 15, Normal_distribution)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)

		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		_, _, array, tree := standardSetup(distanceMatrix, labels)
		_, resulting_tree := neighborJoin(distanceMatrix, labels, array, tree)
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

func Test_Split_Distance(t *testing.T) {
	iterations := 500
	results := make(map[int]float64)
	for i := 0; i < iterations; i++ {

		//GENERATE 2 TREES
		tree1, _, _ := generateTree(5, 15, Normal_distribution)
		tree2, _, _ := generateTree(5, 15, Normal_distribution)

		//CHECK TREES
		results[Split_Distance(tree1[0], tree2[2])]++
	}

	for k, v := range results {
		results[k] = float64(v) / float64(iterations)
	}

	//WE EXPECT 1/15 OF RANDOM 5-tip TREES TO BE TOPOLOGICALLY IDENTICAL
	test := math.Abs(results[0] - float64(1)/float64(15))
	if test > 0.025 {
		t.Errorf("Expect 1/15 good trees, %f", test)
	}

	//CHECK THAT BIGGER TREES DO NOT RANDOMLY MAKE THE SAME TREE
	for i := 0; i < (iterations / 5); i++ {

		//GENERATE 2 TREES
		tree1, _, _ := generateTree(20, 15, Normal_distribution)
		tree2, _, _ := generateTree(20, 15, Normal_distribution)

		//CHECK TREES
		result := Split_Distance(tree1[0], tree2[2])
		if result == 0 {
			t.Errorf("Unlikely scenario - Big trees not expected to randomly be identical")
		}
	}
}
func Test_Canonical_rapid_generate_identical_matrixes_and_split_distance0(t *testing.T) {
	_, labels, distanceMatrix := generateTree(100, 15, Normal_distribution)
	original_labels := make([]string, len(labels))
	copy(original_labels, labels)

	original_dist_mat := make([][]float64, len(distanceMatrix))
	for i := range distanceMatrix {
		original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
		copy(original_dist_mat[i], distanceMatrix[i])
	}

	_, _, array, tree := standardSetup(distanceMatrix, labels)
	_, canon_tree := neighborJoin(distanceMatrix, labels, array, tree)
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

	S, dead_record, array, tree := standardSetup(dist_mat_cpy, labels_cpy)
	_, rapid_tree := rapidJoin(dist_mat_cpy, S, labels_cpy, dead_record, array, tree)
	emptyMatrix2 := make([][]float64, len(labels_cpy))
	for i := range dist_mat_cpy {
		emptyMatrix2[i] = make([]float64, len(labels_cpy))
	}

	resulting_rapid_matrix := createDistanceMatrix(emptyMatrix2, rapid_tree, original_labels)

	cmp_canon_original := compareDistanceMatrixes(original_dist_mat, resulting_canonical_matrix)
	cmp_rapid_original := compareDistanceMatrixes(original_dist_mat, resulting_rapid_matrix)

	if !(cmp_canon_original) {
		t.Errorf(" canon != origninal")
	}
	if !(cmp_rapid_original) {
		t.Errorf(" rapid != original")
	}
	cmp_canon_rapid := compareDistanceMatrixes(resulting_canonical_matrix, resulting_rapid_matrix)
	//this case should not be possible of passed the two other comparisons
	if !(cmp_canon_rapid) {
		t.Errorf("canon != rapid")
	}

	if Split_Distance(canon_tree[0], rapid_tree[0]) != 0 {
		t.Errorf("rapid and canonical tree are not the same")
	}

}

func Test_Compare_runtimes_canonical_against_rapid(t *testing.T) {
	var time_start, time_end int64
	var time_measured int
	NewickFlag = false

	_, labels, distanceMatrix := generateTree(1000, 15, Shifting_Normal_Distribution)
	original_labels := make([]string, len(labels))
	copy(original_labels, labels)

	original_dist_mat := make([][]float64, len(distanceMatrix))
	for i := range distanceMatrix {
		original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
		copy(original_dist_mat[i], distanceMatrix[i])
	}

	_, _, array, tree := standardSetup(distanceMatrix, labels)

	//run rapidJoin and measure the time
	time_start = time.Now().UnixMilli()
	fmt.Printf("###BEGINNING NJ###\n")
	neighborJoin(distanceMatrix, labels, array, tree)
	time_end = time.Now().UnixMilli()
	time_measured = int(time_end - time_start)
	fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured)

	emptyMatrix1 := make([][]float64, len(labels))
	for i := range distanceMatrix {
		emptyMatrix1[i] = make([]float64, len(labels))
	}

	labels_cpy := make([]string, len(original_labels))
	copy(labels_cpy, original_labels)

	dist_mat_cpy := make([][]float64, len(original_dist_mat))
	for i := range original_dist_mat {
		dist_mat_cpy[i] = make([]float64, len(original_dist_mat[i]))
		copy(dist_mat_cpy[i], original_dist_mat[i])
	}

	S, dead_record, array, tree := standardSetup(dist_mat_cpy, labels_cpy)

	//run rapidJoin and measure the time
	time_start = time.Now().UnixMilli()
	fmt.Printf("###BEGINNING RAPIDNJ###\n")
	rapidJoin(dist_mat_cpy, S, labels_cpy, dead_record, array, tree)
	time_end = time.Now().UnixMilli()
	time_measured = int(time_end - time_start)
	fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured)

}

func Test_Make_Time_Taxa_CSV(t *testing.T) {
	taxavalue := 100
	csvFile, err := os.Create("time_plot.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)

	label := []string{"taxa", "rapidnj_shifted", "canonical_shifted", "canonical_norm", "rapid_norm"}
	csvWriter.Write(label)

	for i := 1; i < 15; i++ {
		var time_start, time_end int64
		NewickFlag = false
		fmt.Println()
		fmt.Printf("###TAXASIZE: %d\n", i*taxavalue)

		//make first tree
		_, labels, distanceMatrix := generateTree(i*taxavalue, 15, Shifting_Normal_Distribution)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)
		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		_, _, array, tree := standardSetup(distanceMatrix, labels)

		//make second tree tree
		_, labels2, distanceMatrixUni := generateTree(i*taxavalue, 15, Normal_distribution)
		original_labels2 := make([]string, len(labels2))
		copy(original_labels2, labels2)

		original_dist_mat2 := make([][]float64, len(distanceMatrixUni))
		for i := range distanceMatrixUni {
			original_dist_mat2[i] = make([]float64, len(distanceMatrixUni[i]))
			copy(original_dist_mat2[i], distanceMatrixUni[i])
		}

		_, _, array2, tree2 := standardSetup(distanceMatrixUni, labels2)

		//run CANONICAL and measure the time on Shifting Norm distance matrix
		time_start = time.Now().UnixMilli()
		fmt.Printf("###BEGINNING NJ###\n")
		neighborJoin(distanceMatrix, labels, array, tree)
		time_end = time.Now().UnixMilli()
		time_measured_nj := int(time_end - time_start)
		fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured_nj)

		emptyMatrix1 := make([][]float64, len(labels))
		for i := range distanceMatrix {
			emptyMatrix1[i] = make([]float64, len(labels))
		}

		labels_cpy := make([]string, len(original_labels))
		copy(labels_cpy, original_labels)

		dist_mat_cpy := make([][]float64, len(original_dist_mat))
		for i := range original_dist_mat {
			dist_mat_cpy[i] = make([]float64, len(original_dist_mat[i]))
			copy(dist_mat_cpy[i], original_dist_mat[i])
		}

		S, dead_record, array, tree := standardSetup(dist_mat_cpy, labels_cpy)

		//run rapidJoin and measure the time on Shifting norm
		time_start = time.Now().UnixMilli()
		fmt.Printf("###BEGINNING RAPIDNJ###\n")
		rapidJoin(dist_mat_cpy, S, labels_cpy, dead_record, array, tree)
		time_end = time.Now().UnixMilli()
		time_measured_rapid := int(time_end - time_start)
		fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured_rapid)

		/// #######################  NEW MATRIX CODE STARTING ####################################
		//run canonical and measure the time on standard norm
		time_start = time.Now().UnixMilli()
		fmt.Printf("###BEGINNING NJ###\n")
		neighborJoin(distanceMatrixUni, labels2, array2, tree2)
		time_end = time.Now().UnixMilli()
		time_measured_nj_second := int(time_end - time_start)
		fmt.Printf("### TIME ELAPSED NORMAL DIST: %d ms\n", time_measured_nj_second)

		labels_cpy_2 := make([]string, len(original_labels2))
		copy(labels_cpy_2, original_labels2)

		dist_mat_cpy_2 := make([][]float64, len(original_dist_mat2))
		for i := range original_dist_mat2 {
			dist_mat_cpy_2[i] = make([]float64, len(original_dist_mat2[i]))
			copy(dist_mat_cpy_2[i], original_dist_mat2[i])
		}

		S, dead_record, array, tree = standardSetup(original_dist_mat2, labels_cpy_2)

		//run rapidJoin and measure the time on standard norm distribution
		time_start = time.Now().UnixMilli()
		fmt.Printf("###BEGINNING RAPIDNJ###\n")
		rapidJoin(dist_mat_cpy_2, S, labels_cpy_2, dead_record, array, tree)
		time_end = time.Now().UnixMilli()
		time_measured_rapid_second := int(time_end - time_start)
		fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured_rapid_second)

		row := []string{strconv.Itoa(i * taxavalue), strconv.Itoa(time_measured_rapid), strconv.Itoa(time_measured_nj),
			strconv.Itoa(time_measured_nj_second), strconv.Itoa(time_measured_rapid_second)}
		_ = csvWriter.Write(row)
	}
	csvWriter.Flush()
	csvFile.Close()
}

//#############################################
//helper functions we use in the test framework
//#############################################

func compareDistanceMatrixes(matrix1 [][]float64, matrix2 [][]float64) bool {

	for i, row := range matrix1 {
		for j := range row {
			if math.Abs(matrix1[i][j]-matrix2[i][j]) > 0.1 {

				fmt.Println("first", matrix1[i][j])
				fmt.Println("second", matrix2[i][j])
				return false
			}
		}
	}
	return true
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

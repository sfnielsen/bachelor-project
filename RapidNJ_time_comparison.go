package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	Shifting_Normal_Distribution string = "Sh_norm"
	Normal_distribution          string = "Norm"
	Uniform_distribution         string = "Uniform"
	Cluster_Normal_Distribution  string = "Cluster_norm"
	Spike_Normal_distribution    string = "Spike_norm"
)

func TestRuntimeOfBigTaxas() {
	//declaring some variables to hold times
	var time_start, time_end int64

	fmt.Println("###GENERATING DISTANCE MATRIX")
	time_start = time.Now().UnixMilli()
	_, labels, distanceMatrix := GenerateTree(1000, 1000, Normal_distribution)
	time_end = time.Now().UnixMilli()
	time_generateTree := time_end - time_start
	fmt.Printf("###Done in %d milliseconds\n", time_generateTree)

	S, dead_record, array, treeBanana := standardSetup(distanceMatrix, labels)

	fmt.Println("###BEGIN NEIGHBOR-JOINING")
	time_start = time.Now().UnixMilli()
	a, b := rapidJoin(distanceMatrix, S, labels, dead_record, array, treeBanana, rapidNeighborJoining)
	time_end = time.Now().UnixMilli()
	time_neighborJoin := int(time_end - time_start)
	fmt.Printf("###Done in %d milliseconds\n", time_neighborJoin)

	if a == "" || b == nil {
		fmt.Print("failure :D")
	}
}

func Test_Compare_runtimes_canonical_against_rapid() {
	var time_start, time_end int64
	var time_measured int

	_, labels, distanceMatrix := GenerateTree(1000, 15, Shifting_Normal_Distribution)
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
	rapidJoin(dist_mat_cpy, S, labels_cpy, dead_record, array, tree, rapidNeighborJoining)
	time_end = time.Now().UnixMilli()
	time_measured = int(time_end - time_start)
	fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured)

}

func Test_Make_Time_Taxa_CSV() {
	taxavalue := 100
	csvFile, err := os.Create("time_plot_canonical_vs_rapid.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)

	label := []string{"taxa", "rapidnj", "canonical", "rapidnj_2"}
	csvWriter.Write(label)

	for i := 1; i < 31; i++ {
		var time_start, time_end int64
		fmt.Println()
		fmt.Printf("###TAXASIZE: %d\n", i*taxavalue)

		//make first tree
		_, labels, distanceMatrix := GenerateTree(i*taxavalue, 15, Uniform_distribution)
		original_labels := make([]string, len(labels))
		copy(original_labels, labels)
		original_dist_mat := make([][]float64, len(distanceMatrix))
		for i := range distanceMatrix {
			original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
			copy(original_dist_mat[i], distanceMatrix[i])
		}

		_, _, array, tree := standardSetup(distanceMatrix, labels)

		//make second tree tree
		_, labels2, distanceMatrixUni := GenerateTree(i*taxavalue, 15, Uniform_distribution)
		original_labels2 := make([]string, len(labels2))
		copy(original_labels2, labels2)

		original_dist_mat2 := make([][]float64, len(distanceMatrixUni))
		for i := range distanceMatrixUni {
			original_dist_mat2[i] = make([]float64, len(distanceMatrixUni[i]))
			copy(original_dist_mat2[i], distanceMatrixUni[i])
		}

		//_, _, array2, tree2 := standardSetup(distanceMatrixUni, labels2)

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
		rapidJoin(dist_mat_cpy, S, labels_cpy, dead_record, array, tree, rapidNeighborJoining)
		time_end = time.Now().UnixMilli()
		time_measured_rapid := int(time_end - time_start)
		fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured_rapid)

		/// #######################  NEW MATRIX CODE STARTING ####################################
		//run canonical and measure the time on standard norm
		//time_start = time.Now().UnixMilli()
		//fmt.Printf("###BEGINNING NJ###\n")
		//neighborJoin(distanceMatrixUni, labels2, array2, tree2)
		//time_end = time.Now().UnixMilli()
		//time_measured_nj_second := int(time_end - time_start)
		//fmt.Printf("### TIME ELAPSED NORMAL DIST: %d ms\n", time_measured_nj_second)

		labels_cpy_2 := make([]string, len(original_labels2))
		copy(labels_cpy_2, original_labels2)

		dist_mat_cpy_2 := make([][]float64, len(original_dist_mat2))
		for i := range original_dist_mat2 {
			dist_mat_cpy_2[i] = make([]float64, len(original_dist_mat2[i]))
			copy(dist_mat_cpy_2[i], original_dist_mat2[i])
		}

		S3, dead_record3, array3, tree3 := standardSetup(original_dist_mat2, labels_cpy_2)

		//run rapidJoin and measure the time on standard norm distribution
		time_start = time.Now().UnixMilli()
		fmt.Printf("###BEGINNING RAPIDNJ###\n")
		rapidJoin(dist_mat_cpy_2, S3, labels_cpy_2, dead_record3, array3, tree3, rapidNeighborJoining)
		time_end = time.Now().UnixMilli()
		time_measured_rapid_second := int(time_end - time_start)
		fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured_rapid_second)

		row := []string{strconv.Itoa(i * taxavalue), strconv.Itoa(time_measured_rapid), strconv.Itoa(time_measured_nj),
			strconv.Itoa(time_measured_rapid_second)}
		_ = csvWriter.Write(row)
	}
	csvWriter.Flush()
	csvFile.Close()
}

func test_all_trees_on_rapidNj() {
	taxavalue := 100
	csvFile, err := os.Create("allTrees_timetest.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)

	treeTypes := []string{Shifting_Normal_Distribution, Normal_distribution, Uniform_distribution,
		Cluster_Normal_Distribution, Spike_Normal_distribution}

	csvWriter.Write(treeTypes)
	NewickFlag = false

	for i := 1; i < 31; i++ {
		var time_start, time_end int64
		fmt.Println()
		fmt.Printf("###TAXASIZE: %d\n", i*taxavalue)

		row := make([]string, len(treeTypes))
		row = append(row, strconv.Itoa(i*taxavalue))
		for _, treeType := range treeTypes {
			array, tree, distanceMatrix, labels, S, dead_record := setupDistanceMatrixForTimeTaking(i, taxavalue, treeType)

			//run rapidJoin and measure the time on Shifting norm
			fmt.Printf(treeType)
			fmt.Printf("###BEGINNING RAPIDNJ###\n")
			time_start = time.Now().UnixMilli()

			rapidJoin(distanceMatrix, S, labels, dead_record, array, tree, rapidNeighborJoining)
			time_end = time.Now().UnixMilli()
			time_measured_rapid := int(time_end - time_start)
			fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured_rapid)
			row = append(row, strconv.Itoa(time_measured_rapid))
		}
		_ = csvWriter.Write(row)
		csvWriter.Flush()

	}
	csvFile.Close()
}


func compare_runtime_on_umax_vs_normal_rapidnj(t *testing.T) {
	_, labels, distanceMatrix := GenerateTree(100, 15, Normal_distribution)
	original_labels := make([]string, len(labels))
	copy(original_labels, labels)

	original_dist_mat := make([][]float64, len(distanceMatrix))
	for i := range distanceMatrix {
		original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
		copy(original_dist_mat[i], distanceMatrix[i])
	}

	_, _, array, tree := standardSetup(distanceMatrix, labels)
	_, canon_tree := rapidJoin(original_dist_mat, S, original_labels, dead_record, array, tree, rapidNeighborJoining)
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
	_, rapid_tree := rapidJoin(dist_mat_cpy, S, labels_cpy, dead_record, array, tree, rapidNeighborJoining)
	emptyMatrix2 := make([][]float64, len(labels_cpy))
	for i := range dist_mat_cpy {
		emptyMatrix2[i] = make([]float64, len(labels_cpy))
	}
}




//helper methods
func standardSetup(D [][]float64, labels []string) ([][]Tuple, map[int]int, Tree, Tree) {
	S := initSmatrix(D)
	deadRecords := initDeadRecords(D)
	var tree Tree
	var label_tree Tree = generateTreeForRapidNJ(labels)

	tree = append(tree, label_tree...)

	return S, deadRecords, label_tree, tree
}

func setupDistanceMatrixForTimeTaking(i int, taxavalue int, treeType string) (Tree, Tree, [][]float64,
	[]string, [][]Tuple, map[int]int) {
	_, labels, distanceMatrix := GenerateTree(i*taxavalue, 15, treeType)
	original_labels := make([]string, len(labels))
	copy(original_labels, labels)
	original_dist_mat := make([][]float64, len(distanceMatrix))
	for i := range distanceMatrix {
		original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
		copy(original_dist_mat[i], distanceMatrix[i])
	}

	S, dead_records, array, tree := standardSetup(distanceMatrix, labels)

	return array, tree, distanceMatrix, labels, S, dead_records

}

func compare_U_max_sorting(){

}

func main() {
	Test_Make_Time_Taxa_CSV()
}

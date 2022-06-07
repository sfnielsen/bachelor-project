package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
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
	seed := time.Now().UTC().UnixNano()
	time_start = time.Now().UnixMilli()
	_, labels, distanceMatrix := GenerateTree(1000, 1000, Normal_distribution, seed)
	time_end = time.Now().UnixMilli()
	time_generateTree := time_end - time_start
	fmt.Printf("###Done in %d milliseconds\n", time_generateTree)

	fmt.Println("###BEGIN NEIGHBOR-JOINING")
	time_start = time.Now().UnixMilli()
	a, b := rapidNeighbourJoin(distanceMatrix, labels, rapidNeighborJoining)
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
	seed := time.Now().UTC().UnixNano()
	_, labels, distanceMatrix := GenerateTree(1000, 15, Normal_distribution, seed)
	original_labels := make([]string, len(labels))
	copy(original_labels, labels)

	original_dist_mat := make([][]float64, len(distanceMatrix))
	for i := range distanceMatrix {
		original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
		copy(original_dist_mat[i], distanceMatrix[i])
	}

	//_, _, array, tree := standardSetup(distanceMatrix, labels)

	//run rapidJoin and measure the time
	time_start = time.Now().UnixMilli()
	fmt.Printf("###BEGINNING NJ###\n")
	neighborJoin(distanceMatrix, labels)
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
	//run rapidJoin and measure the time
	time_start = time.Now().UnixMilli()
	fmt.Printf("###BEGINNING RAPIDNJ###\n")
	rapidNeighbourJoin(dist_mat_cpy, labels_cpy, rapidNeighborJoining)
	time_end = time.Now().UnixMilli()
	time_measured = int(time_end - time_start)
	fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured)

}
func Test_make_rapid_u_updates_CSV() {
	taxavalue := 100
	csvFile, err := os.Create("radix_time.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)

	label := []string{"taxa", "rapidnj", "rapidnj_error", "mean"}
	csvWriter.Write(label)

	for i := 1; i < 63; i++ {
		highest_rapidnj, lowest_rapidnj := 0, 9999999999999999
		mean_rapidnj := 0

		iterations := 10
		seed := int64(2345)
		fmt.Printf("###TAXASIZE: %d\n", i*taxavalue)
		for j := 0; j < iterations; j++ {
			seed++
			var time_start, time_end int64

			//make first tree
			_, labels, distanceMatrix := GenerateTree(i*taxavalue, 100, Normal_distribution, seed)
			original_labels := make([]string, len(labels))
			copy(original_labels, labels)
			original_dist_mat := make([][]float64, len(distanceMatrix))
			for i := range distanceMatrix {
				original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
				copy(original_dist_mat[i], distanceMatrix[i])
			}

			labels_cpy := make([]string, len(original_labels))
			copy(labels_cpy, original_labels)

			dist_mat_cpy := make([][]float64, len(original_dist_mat))
			for i := range original_dist_mat {
				dist_mat_cpy[i] = make([]float64, len(original_dist_mat[i]))
				copy(dist_mat_cpy[i], original_dist_mat[i])
			}

			//run rapidJoin and measure the time on Shifting norm
			time_start = time.Now().UnixMilli()
			neighborJoin(dist_mat_cpy, labels_cpy)
			time_end = time.Now().UnixMilli()
			time_measured_rapid := int(time_end - time_start)

			mean_rapidnj += time_measured_rapid
			//finding if time was extrema
			if time_measured_rapid > highest_rapidnj {
				highest_rapidnj = time_measured_rapid
			}
			if time_measured_rapid < int(lowest_rapidnj) {
				lowest_rapidnj = (time_measured_rapid)
			}

		}
		fmt.Printf("### TIME ELAPSED mean: %d ms\n", mean_rapidnj)
		fmt.Println(highest_rapidnj, lowest_rapidnj)
		row := []string{strconv.Itoa(i * taxavalue),
			strconv.Itoa((highest_rapidnj-int(lowest_rapidnj))/2 + lowest_rapidnj),
			fmt.Sprintf("%v", (math.Log(float64(highest_rapidnj))-math.Log(float64(lowest_rapidnj)))/2),
			fmt.Sprintf("%v", mean_rapidnj/iterations)}
		_ = csvWriter.Write(row)
		csvWriter.Flush()
	}
	csvFile.Close()
}

func test_record_all_points(distribution string, filename string, maxlength int) {
	taxavalue := 100
	csvFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)

	for i := 1; i < 61; i++ {
		highest_rapidnj, lowest_rapidnj := 0, 9999999999999999
		mean_rapidnj := 0

		iterations := 10
		seed := int64(2345)
		fmt.Printf("###TAXASIZE: %d\n", i*taxavalue)
		times := []string{strconv.Itoa(i * taxavalue)}
		for j := 0; j < iterations; j++ {
			seed++
			var time_start, time_end int64

			//make first tree
			_, labels, distanceMatrix := GenerateTree(i*taxavalue, maxlength, distribution, seed)
			original_labels := make([]string, len(labels))
			copy(original_labels, labels)
			original_dist_mat := make([][]float64, len(distanceMatrix))
			for i := range distanceMatrix {
				original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
				copy(original_dist_mat[i], distanceMatrix[i])
			}

			labels_cpy := make([]string, len(original_labels))
			copy(labels_cpy, original_labels)

			dist_mat_cpy := make([][]float64, len(original_dist_mat))
			for i := range original_dist_mat {
				dist_mat_cpy[i] = make([]float64, len(original_dist_mat[i]))
				copy(dist_mat_cpy[i], original_dist_mat[i])
			}

			//run rapidJoin and measure the time on Shifting norm
			time_start = time.Now().UnixMilli()

			rapidNeighbourJoin(dist_mat_cpy, labels_cpy, rapidNeighborJoining)
			time_end = time.Now().UnixMilli()
			time_measured_rapid := int(time_end - time_start)

			times = append(times, strconv.Itoa(time_measured_rapid))

		}
		fmt.Printf("### TIME ELAPSED mean: %d ms\n", mean_rapidnj)
		fmt.Println(highest_rapidnj, lowest_rapidnj)
		_ = csvWriter.Write(times)
		csvWriter.Flush()
	}
	csvFile.Close()
}

func Test_Make_Time_Taxa_CSV() {
	taxavalue := 100
	csvFile, err := os.Create("time_plot_canonical_vs_rapid_from_u_update.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)

	label := []string{"taxa", "rapidnj", "canonical", "rapidnj_error", "canonical_error"}
	csvWriter.Write(label)

	for i := 1; i < 45; i++ {
		highest_canonical, lowest_canonical, highest_rapidnj, lowest_rapidnj := 0, 9999999999999999, 0, 999999999999999
		mean_rapidnj, mean_canonical := 0, 0

		iterations := 10
		seed := int64(29638)
		for j := 0; j < iterations; j++ {
			seed++
			var time_start, time_end int64
			fmt.Println()
			fmt.Printf("###TAXASIZE: %d\n", i*taxavalue)

			//make first tree
			_, labels, distanceMatrix := GenerateTree(i*taxavalue, 15, Uniform_distribution, seed)
			original_labels := make([]string, len(labels))
			copy(original_labels, labels)
			original_dist_mat := make([][]float64, len(distanceMatrix))
			for i := range distanceMatrix {
				original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
				copy(original_dist_mat[i], distanceMatrix[i])
			}

			//_, _, array, tree := standardSetup(distanceMatrix, labels)

			//make second tree tree
			_, labels2, distanceMatrixUni := GenerateTree(i*taxavalue, 15, Uniform_distribution, seed)
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
			neighborJoin(distanceMatrix, labels)
			time_end = time.Now().UnixMilli()
			time_measured_nj := int(time_end - time_start)
			fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured_nj)

			//finding if time was extrema
			if time_measured_nj > highest_canonical {
				highest_canonical = time_measured_nj
			}
			if time_measured_nj < int(lowest_canonical) {
				lowest_canonical = time_measured_nj
			}

			mean_canonical += time_measured_nj

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

			//run rapidJoin and measure the time on Shifting norm
			time_start = time.Now().UnixMilli()
			fmt.Printf("###BEGINNING RAPIDNJ###\n")
			rapidNeighbourJoin(dist_mat_cpy, labels_cpy, rapidNeighborJoining)
			time_end = time.Now().UnixMilli()
			time_measured_rapid := int(time_end - time_start)
			fmt.Printf("### TIME ELAPSED: %d ms\n", time_measured_rapid)

			mean_rapidnj += time_measured_rapid
			//finding if time was extrema
			if time_measured_rapid > highest_rapidnj {
				highest_rapidnj = time_measured_rapid
			}
			if time_measured_rapid < int(lowest_rapidnj) {
				lowest_rapidnj = (time_measured_rapid)
			}

		}
		fmt.Println(highest_canonical, lowest_canonical, highest_rapidnj, lowest_rapidnj)
		row := []string{strconv.Itoa(i * taxavalue),
			strconv.Itoa((highest_rapidnj-int(lowest_rapidnj))/2 + lowest_rapidnj),
			strconv.Itoa((highest_canonical-int(lowest_canonical))/2 + lowest_canonical),
			fmt.Sprintf("%v", (math.Log(float64(highest_rapidnj))-math.Log(float64(lowest_rapidnj)))/2),
			fmt.Sprintf("%v", (math.Log(float64(highest_canonical))-math.Log(float64(lowest_canonical)))/2)}

		_ = csvWriter.Write(row)
		csvWriter.Flush()
	}
	csvFile.Close()
}

func test_all_trees_on_rapidNj() {
	itrsize := 10
	taxavalue := 100
	csvFile, err := os.Create("allTrees_timetest.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)

	treeTypes := []string{Normal_distribution, Cluster_Normal_Distribution, Spike_Normal_distribution}
	errs := []string{"norm_err", "cluster_err", "spike_err"}

	labels := append(treeTypes, errs...)
	labels = append([]string{"taxa"}, labels...)

	csvWriter.Write(labels)
	NewickFlag = false

	for i := 1; i < 31; i++ {
		NewickFlag = false
		var time_start, time_end int64
		fmt.Println()
		fmt.Printf("###TAXASIZE: %d\n", 2900+i*taxavalue)

		row := make([]string, 0)
		row = append(row, strconv.Itoa(int((float64(2900 + i*taxavalue)))))
		errors_row := make([]string, 0)

		seed := int64(7102)
		for _, treeType := range treeTypes {

			mean_rapidnj, highest_rapidnj, lowest_rapidnj := 0, 0, 99999999999999999

			fmt.Println(treeType)
			fmt.Printf("###BEGINNING RAPIDNJ###\n")
			for j := 0; j < itrsize; j++ {
				seed++
				_, labels, distanceMatrix := GenerateTree(2900+i*taxavalue, 100, treeType, seed)

				time_start = time.Now().UnixMilli()
				rapidNeighbourJoin(distanceMatrix, labels, rapidNeighborJoining)
				time_end = time.Now().UnixMilli()
				time_measured_rapid := int(time_end - time_start)

				mean_rapidnj += time_measured_rapid
				//finding if time was extrema
				if time_measured_rapid > highest_rapidnj {
					highest_rapidnj = time_measured_rapid
				}
				if time_measured_rapid < int(lowest_rapidnj) {
					lowest_rapidnj = (time_measured_rapid)
				}
			}
			fmt.Printf("### TIME ELAPSED: %d ms\n", mean_rapidnj/itrsize)
			row = append(row, strconv.Itoa((highest_rapidnj-lowest_rapidnj)/2+lowest_rapidnj))
			errors_row = append(errors_row, fmt.Sprintf("%v", (math.Log(float64(highest_rapidnj))-math.Log(float64(lowest_rapidnj)))/2))

		}
		row = append(row, errors_row...)
		_ = csvWriter.Write(row)
		csvWriter.Flush()

	}
	csvFile.Close()
}

func compare_runtime_on_umax_vs_normal_rapidnj() {
	seed := time.Now().UTC().UnixNano()
	_, labels, distanceMatrix := GenerateTree(1500, 15, Cluster_Normal_Distribution, seed)
	original_labels := make([]string, len(labels))
	copy(original_labels, labels)

	original_dist_mat := make([][]float64, len(distanceMatrix))
	for i := range distanceMatrix {
		original_dist_mat[i] = make([]float64, len(distanceMatrix[i]))
		copy(original_dist_mat[i], distanceMatrix[i])
	}

	labels_cpy := make([]string, len(original_labels))
	copy(labels_cpy, original_labels)

	var time_start float64
	var time_end float64
	//STANDARD RAPID JOINING

	time_start = float64(time.Now().UnixMilli())
	rapidNeighbourJoin(distanceMatrix, labels, rapidNeighborJoining)
	time_end = float64(time.Now().UnixMilli())
	standard_rapid_time := time_end - time_start
	fmt.Println(standard_rapid_time)

	dist_mat_cpy := make([][]float64, len(original_dist_mat))
	for i := range original_dist_mat {
		dist_mat_cpy[i] = make([]float64, len(original_dist_mat[i]))
		copy(dist_mat_cpy[i], original_dist_mat[i])
	}

	//U SORTED RAPID JOINING
	time_start = float64(time.Now().UnixMilli())
	rapidNeighbourJoin(dist_mat_cpy, labels_cpy, rapidNeighborJoining)
	time_end = float64(time.Now().UnixMilli())
	U_max_heuristic_time := time_end - time_start
	fmt.Println(U_max_heuristic_time)
}

// #####################################################################################################################################
// #####################################################################################################################################
//helper methods

/*
func standardSetup(D [][]float64, labels []string) ([][]Tuple, map[int]int, Tree, Tree) {
	S := initSmatrix(D)
	deadRecords := initLiveRecords(D)
	var tree Tree
	var label_tree Tree = generateTreeForRapidNJ(labels)

	tree = append(tree, label_tree...)

	return S, deadRecords, label_tree, tree
}*/

/*
func setupDistanceMatrixForTimeTaking(i int, taxavalue int, treeType string) (Tree, Tree, [][]float64,
	[]string, [][]Tuple, map[int]int) {
	seed := time.Now().UTC().UnixNano()
	tree, labels, distanceMatrix := GenerateTree(i*taxavalue, 15, treeType, seed)

	S, live_records, array, _ := standardSetup(distanceMatrix, labels)

	return array, tree, distanceMatrix, labels, S, live_records

}
*/
func compare_U_max_sorting() {

}

func main() {
	Test_make_rapid_u_updates_CSV()
}

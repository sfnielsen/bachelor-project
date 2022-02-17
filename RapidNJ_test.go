package main

import (
	"testing"
)

func standardSetup(D [][]float64) ([][]Tuple, map[int]int) {
	S := initSmatrix(D)
	deadRecords := initDeadRecords(D)
	return S, deadRecords
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
	S, deadRecords := standardSetup(D)

	newick_result := neighborJoin(D, S, labels, deadRecords)
	if newick_result != "((B:4.000000,A:13.000000):2.000000,(C:4.000000,D:10.000000):2.000000);" {
		t.Errorf("hehehe")
	}

}

func Test4Taxa_made_up_numbers(t *testing.T) {
	labels := []string{
		"A", "B", "C", "D",
	}
	D := [][]float64{
		{0, 11, 18, 24},
		{11, 0, 12, 18},
		{18, 12, 0, 14},
		{24, 18, 14, 0},
	}
	S, deadRecords := standardSetup(D)

	newick_result := neighborJoin(D, S, labels, deadRecords)
	if newick_result != "((B:2.500000,A:8.500000):2.750000,(C:4.000000,D:10.000000):2.750000);" {
		t.Errorf("hehehe")
	}

}

func Test8Taxa_madeUpNumbers_shouldBeChangedLater(t *testing.T) {

	labels := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	D := [][]float64{
		{0, 5, 68, 57, 127, 27, 28, 33},
		{5, 0, 58, 47, 117, 8, 52, 57},
		{68, 58, 0, 35, 69, 35, 87, 92},
		{57, 47, 35, 0, 94, 44, 79, 84},
		{127, 117, 69, 94, 0, 144, 149, 154},
		{27, 8, 35, 44, 144, 0, 27, 54},
		{28, 52, 87, 79, 149, 27, 0, 13},
		{33, 57, 92, 84, 154, 54, 13, 0},
	}

	S, deadRecords := standardSetup(D)

	newick_result := neighborJoin(D, S, labels, deadRecords)
	if newick_result != "(B:5.765625,((((G:1.250000,H:11.750000):23.208333,A:0.791667):4.718750,F:29.781250):1.281250,((C:0.333333,E:68.666667):18.200000,D:11.800000):2.968750):5.765625);" {
		t.Errorf("hehehe")
	}
}

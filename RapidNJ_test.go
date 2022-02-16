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

package main

import (
	"math"
	"sort"
)

type S_Search_Strategy func(u []float64, D [][]float64, S [][]Tuple, live_records map[int]int) (int, int)

func rapidNeighborJoining_U_sorted(u []float64, D [][]float64, S [][]Tuple, live_records map[int]int) (int, int) {
	max_u := MaxIntSlice(u)
	q_min := math.MaxFloat64
	cur_i, cur_j := -1, -1

	//begin u-max ideaimpl
	u_order := make([]*U_Tuple, 0)
	for i, v := range u {
		new_tuple := new(U_Tuple)
		new_tuple.index_in_d = i
		new_tuple.value = v
		u_order = append(u_order, new_tuple)
	}
	sort.Slice(u_order, func(a, b int) bool {
		return (u_order[a].value > u_order[b].value)
	})

	for _, v := range u_order {
		for c := range S[v.index_in_d] {
			s := S[v.index_in_d][c]
			c_to_cD, ok := live_records[s.index_j]
			//check if dead record
			if !ok {
				continue
			}
			// case where i == j
			if v.index_in_d == c_to_cD {
				continue
			}
			if s.value-u[v.index_in_d]-max_u > q_min {
				break
			}
			q := s.value - u[v.index_in_d] - u[c_to_cD]
			if q < q_min {
				cur_i = v.index_in_d
				cur_j = c_to_cD
				q_min = q
			}

		}
	}

	return cur_i, cur_j
}

func rapidNeighborJoining(u []float64, D [][]float64, S [][]Tuple, live_records map[int]int) (int, int) {
	max_u := MaxIntSlice(u)
	q_min := math.MaxFloat64
	cur_i, cur_j := -1, -1

	for r, row := range S {

		for c := range row {
			s := S[r][c]
			c_to_cD, ok := live_records[s.index_j]
			//check if dead record
			if !ok {
				continue
			}
			// case where i == j
			if r == c_to_cD {
				continue
			}
			if s.value-u[r]-max_u > q_min {
				break
			}
			q := s.value - u[r] - u[c_to_cD]
			if q < q_min {
				cur_i = r
				cur_j = c_to_cD
				q_min = q
			}
		}
	}

	return cur_i, cur_j
}

package main

import (
	"math"
	"sort"
	"time"
)

type S_Search_Strategy func(u []float64, D [][]float64, S [][]Tuple, live_records map[int]int) (int, int)

func rapidNeighborJoining_U_sorted(u []float64, D [][]float64, S [][]Tuple, live_records map[int]int, row_for_new_it int) (int, int) {
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

//####################################
//vars for performing tests on rapid NJ
//####################################
var total_lookups, total_updates int
var old_i int = -1
var column_depth map[int]int = make(map[int]int)

var extra_cost int
var last_q_min float64 = math.MaxFloat64

var start, end, total int64

var lookup_updates_count = false

const taxa_lookup_update = 400

var lookup_matrix = make([][]int, taxa_lookup_update)
var update_matrix = make([][]int, taxa_lookup_update)

//####################################
//####################################
//####################################

func mapLookupInSearchHeurisic(blub int, live_records map[int]int) (int, bool) {
	c_to_cD, ok := live_records[blub]
	return c_to_cD, ok
}

func finQMin(u []float64, D [][]float64, S [][]Tuple, live_records map[int]int) (float64, int, int) {
	q_min := math.MaxFloat64
	cur_i, cur_j := -1, -1

	for r, row := range S {
		for c := range row {
			if c > 1 {
				break
			}

			if c == 0 {
				continue
			}

			s := S[r][c]

			time_start = time.Now().UnixNano()

			c_to_cD, ok := mapLookupInSearchHeurisic(s.index_j, live_records)
			time_end = time.Now().UnixNano()
			lookupTime += (int(time_end) - int(time_start))

			//check if dead record
			if !ok {
				continue
			}

			q := s.value - u[r] - u[c_to_cD]

			if q < q_min {
				cur_i = r
				cur_j = c_to_cD
				q_min = q
			}
		}
	}
	return q_min, cur_i, cur_j
}

func rapidNeighborJoining(u []float64, D [][]float64, S [][]Tuple, live_records map[int]int) (int, int) {
	max_u := MaxIntSlice(u)
	q_min, cur_i, cur_j := finQMin(u, D, S, live_records)
	for r, row := range S {
		for c := range row {

			if c == 0 {
				continue
			}

			s := S[r][c]

			if s.value-u[r]-max_u > q_min {
				break
			}

			if lookup_updates_count && len(D) == taxa_lookup_update {
				lookup_matrix[r][c]++
			}
			time_start = time.Now().UnixNano()

			c_to_cD, ok := mapLookupInSearchHeurisic(s.index_j, live_records)
			time_end = time.Now().UnixNano()
			lookupTime += (int(time_end) - int(time_start))

			//check if dead record
			if !ok {
				continue
			}

			q := s.value - u[r] - u[c_to_cD]

			if q < q_min {
				if lookup_updates_count && len(D) == taxa_lookup_update {
					update_matrix[r][c]++
				}
				cur_i = r
				cur_j = c_to_cD
				q_min = q

			}
		}
	}
	return cur_i, cur_j
}

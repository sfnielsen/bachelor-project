package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func Parse_text(file string) ([][]float64, []string) {

	f, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	var wordno = -1
	var word string
	var taxa int
	var value float64
	var row, phylip_col int
	var D = make([][]float64, 0)
	var labels = make([]string, 0)
	for scanner.Scan() {

		word = scanner.Text()
		if word == "" {
			continue
		}

		//first word should be amount of taxa, which dictates dimensions of D
		if wordno == -1 {
			taxa, err = strconv.Atoi(scanner.Text())
			if err != nil {
				panic(err)
			}

			for r := 0; r < taxa; r++ {
				D = append(D, make([]float64, taxa))
			}

			wordno++
			continue
		}

		//check if label, if not then it is value to be stored in D
		phylip_col = wordno % (taxa + 1)
		if phylip_col == 0 {

			row = wordno / taxa
			labels = append(labels, word)

		} else {

			//default case - we insert a value to D
			value, err = strconv.ParseFloat(scanner.Text(), 64)
			if err != nil {
				panic(err)
			}
			col := phylip_col - 1
			D[row][col] = value

		}
		wordno++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return D, labels
}

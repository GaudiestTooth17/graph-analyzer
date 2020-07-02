package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gitlab.com/hcmi/graph-analyzer/adjmat"
)

func readAdjacencyMatrix(fileName string) [][]uint16 {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)

	//set up adjMatrix
	line, err := fileReader.ReadString('\n')
	row, err := lineToSlice(line)
	if err != nil {
		fmt.Println(line)
		fmt.Println("Cannot be read as a row of the matrix.")
		panic("")
	}
	adjMatrix := make([][]uint16, len(row))
	for i := 0; i < len(row); i++ {
		adjMatrix[i] = make([]uint16, len(row))
	}
	adjMatrix[0] = row

	//populate adjMatrix
	for i := 1; i < len(row); i++ {
		line, err = fileReader.ReadString('\n')
		if err != nil {
			panic("Could not read line " + strconv.Itoa(i))
		}
		adjMatrix[i], err = lineToSlice(line)
		if err != nil {
			fmt.Printf("\nInvalid format on line %d (printed above).\n", i)
			panic("")
		}
	}

	return adjMatrix
}

func lineToSlice(line string) ([]uint16, error) {
	slice := make([]uint16, 0)
	numbers := strings.Fields(line)

	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Println("An error occurred while reading to the file.")
			return nil, err
		}
		slice = append(slice, uint16(num))
	}

	return slice, nil
}

func writeMatrix(fileName string, matrix *adjmat.AdjacencyMatrix) {
	rows, cols := (*matrix).Dims()
	outFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Fprintf(outFile, "%d ", int((*matrix).At(i, j)))
		}
		fmt.Fprintln(outFile, "")
	}
}

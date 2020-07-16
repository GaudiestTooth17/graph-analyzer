package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
		panic("The above line cannot be read as a row of the matrix.")
	}
	adjMatrix := makeAdjacencyMatrix(len(row))
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

func readAdjacencyList(fileName string) [][]uint16 {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)

	//set up adjMatrix
	line, err := fileReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	n, err := strconv.Atoi(line[:len(line)-1])
	if err != nil {
		panic(err)
	}
	adjMatrix := makeAdjacencyMatrix(n)

	//populate adjMatrix
	for true {
		line, err = fileReader.ReadString('\n')
		// This is the exit condition
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil && !errors.Is(err, io.EOF) {
			panic(err)
		}
		i, j := lineToCoordinate(line)
		adjMatrix[i][j] = 1
	}
	return adjMatrix
}

// makeAdjacencyMatrix creates an empty adjacency matrix of size n x n
func makeAdjacencyMatrix(n int) [][]uint16 {
	adjMatrix := make([][]uint16, n)
	for i := 0; i < n; i++ {
		adjMatrix[i] = make([]uint16, n)
	}
	return adjMatrix
}

func lineToCoordinate(line string) (i, j int) {
	coordinates := strings.Fields(line)
	if len(coordinates) != 2 {
		panic("Error: " + line)
	}
	i, err := strconv.Atoi(coordinates[0])
	if err != nil {
		panic(err)
	}
	j, err = strconv.Atoi(coordinates[1])
	if err != nil {
		panic(err)
	}
	return i, j
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

// writeMatrix writes an adjacency matrix to a file in matrix form
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

// writeList writes an adjacency matrix to a file in list form
func writeAdjacencyList(fileName string, matrix *adjmat.AdjacencyMatrix) {
	n, _ := (*matrix).Dims()
	outFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// write size of matrix, then iterate over just the necessary parts of the matrix
	fmt.Fprintf(outFile, "%d\n", n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if (*matrix).At(i, j) == 0 {
				continue
			}
			fmt.Fprintf(outFile, "%d %d\n", i, j)
		}
	}
}

func writeMedianDistanceCSV(fileName string, medianDistances []uint16) {
	outFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	for i := 0; i < len(medianDistances); i++ {
		fmt.Fprintf(outFile, "%d, %d\n", i, medianDistances[i])
	}
}

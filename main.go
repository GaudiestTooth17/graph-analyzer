package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/hcmi/graph-analyzer/adjmat"
	"gonum.org/v1/gonum/blas/blas32"
	blas_netlib "gonum.org/v1/netlib/blas/netlib"
)

type void struct{}

const modeMessage = "Modes are comp(onent), dia(meter), dia(meter) fr(om) a list of nodes, dist(ance) distribution."

// ex ./graph-analyzer diafr data/monday.txt data/monday_giant_component.txt
func main() {
	// TODO add help function
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [mode] [file]\n", os.Args[0])
		fmt.Println(modeMessage)
		return
	}
	taskSelection := os.Args[1]
	timeStart := time.Now()
	fmt.Printf("Started %s at %v.\n", taskSelection, timeStart)

	// use the c based library
	blas32.Use(blas_netlib.Implementation{})

	// do what the user said to
	t, matrix := newTask()
	t.run(matrix)

	fmt.Printf("Done with task %s (%v).\n", taskSelection, time.Now().Sub(timeStart))
}

func toAdjacencyMatrix(input [][]uint16) *adjmat.AdjacencyMatrix {
	nodes := len(input)
	data := make([]float32, nodes*nodes)
	for i := 0; i < nodes; i++ {
		for k := 0; k < nodes; k++ {
			data[i*nodes+k] = float32(input[i][k])
		}
	}
	return adjmat.New(nodes, data)
}

type taskFunction func([][]uint16)

type task struct {
	run taskFunction
}

func comp(adjacencyMatrix [][]uint16) {
	components := findComponents(adjacencyMatrix)

	// find largest component
	largestSize := -1
	largestIndex := -1
	for i := 0; i < len(components); i++ {
		if len(components[i]) > largestSize {
			largestSize = len(components[i])
			largestIndex = i
		}
	}

	// create subgraph from largest component
	var largestSubGraph [][]uint16
	if len(components) > 1 {
		largestSubGraph = createGraphFromComponent(components[largestIndex], adjacencyMatrix)
	} else {
		largestSubGraph = adjacencyMatrix
	}

	writeMatrix("lc.txt", toAdjacencyMatrix(largestSubGraph))
}

func dia(adjacencyMatrix [][]uint16) {
	diameter := calculateDiameter(toAdjacencyMatrix(adjacencyMatrix))
	fmt.Printf("The diameter %d.\n", diameter)
}

func dist(adjacencyMatrix [][]uint16) {
	distanceDistribution := calculateDistDistribution(toAdjacencyMatrix(adjacencyMatrix))
	fmt.Println("distance,frequency")
	// display freqency distribution
	for distance, frequency := range distanceDistribution {
		fmt.Printf("%d,%d\n", distance, frequency)
	}
	savePlot(getFileName(), distanceDistribution)
}

func getFileName() string {
	parts := strings.Split(os.Args[2], "/")
	return parts[len(parts)-1]
}

//newTask returns a task and the matrix to run the task on
func newTask() (task, [][]uint16) {
	taskType := os.Args[1]
	fileName := os.Args[2]

	// read graph from HDD
	timeGraph := time.Now()
	graph := readAdjacencyMatrix(fileName)
	fmt.Printf("Read %dx%d adjacency matrix (%v).\n", len(graph), len(graph[0]), time.Now().Sub(timeGraph))

	var fn taskFunction
	if taskType == "comp" {
		fn = comp
	} else if taskType == "dia" {
		fn = dia
	} else if taskType == "diafr" {
		allowedNodes := readAllowedNodesFile(os.Args[3])
		graph = createGraphFromComponent(allowedNodes, graph)
		fn = dia
	} else if taskType == "dist" {
		fn = dist
	} else {
		fn = invalidSelection
	}
	return task{run: fn}, graph
}

func invalidSelection(_ [][]uint16) {
	fmt.Println(modeMessage)
}

func readAllowedNodesFile(fileName string) map[int]void {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)
	scanner := bufio.NewScanner(fileReader)

	allowedNodes := make(map[int]void)
	for scanner.Scan() {
		node, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		allowedNodes[node] = void{}
	}

	return allowedNodes
}

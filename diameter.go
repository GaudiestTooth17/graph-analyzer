package main

import (
	"fmt"
	"math"
	"time"

	"gitlab.com/hcmi/graph-analyzer/adjmat"
	"gitlab.com/hcmi/graph-analyzer/trackingset"
)

func calculateDiameter(adjacencyMatrix *adjmat.AdjacencyMatrix) int {
	iterations := 0

	// resultMatrix starts off representing paths of length 1
	tSet := initializeTrackingSet(adjacencyMatrix)
	resultMatrix := adjmat.CopyOf(adjacencyMatrix)

	previousElems := uint32(math.MaxUint32)
	// each iteration resultMatrix is changed to represent paths of length 1 greater than before
	// once all the nodes are connected, the loop will execute one more time
	// the value in previousElems won't change and the loop will terminate
	for previousElems != tSet.Size() {
		timeStart := time.Now()
		previousElems = tSet.Size()
		resultMatrix = adjmat.Multiply(resultMatrix, adjacencyMatrix) // the hangup is here :(
		tSet.AddAllNonzero(resultMatrix)                              // takes essentially no time
		iterations++
		fmt.Printf("Completed loop %d (%v).\n", iterations, time.Now().Sub(timeStart))
	}

	return iterations
}

func initializeTrackingSet(adjacencyMatrix *adjmat.AdjacencyMatrix) trackingset.Set {
	rows, cols := (*adjacencyMatrix).Dims()
	tSet := trackingset.NewSet(rows, cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if i == j || (*adjacencyMatrix).At(i, j) != 0 {
				tSet.Add(i, j)
			}
		}
	}
	return tSet
}

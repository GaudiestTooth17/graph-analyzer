package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"gitlab.com/hcmi/graph-analyzer/adjmat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func calculateDistDistribution(adjacencyMatrix *adjmat.AdjacencyMatrix) map[uint16]uint16 {
	distMatrix := makeDistMatrix(adjacencyMatrix)
	distDistribution := make(map[uint16]uint16)

	var distance uint16
	// the loop looks weird because the matrix is symmetric
	for i := 0; i < len(distMatrix); i++ {
		for j := i + 1; j < len(distMatrix); j++ {
			distance = distMatrix[i][j]
			if distance == 0 {
				continue
			}
			if _, ok := distDistribution[distance]; !ok {
				distDistribution[distance] = 0
			}
			distDistribution[distance]++
		}
	}

	return distDistribution
}

func makeDistMatrix(adjacencyMatrix *adjmat.AdjacencyMatrix) [][]uint16 {
	//initialize minDistMatrix
	rows, _ := adjacencyMatrix.Dims()
	minDistMatrix := make([][]uint16, rows)
	for i := 0; i < rows; i++ {
		minDistMatrix[i] = make([]uint16, rows)
	}
	resultMatrix := adjmat.CopyOf(adjacencyMatrix)
	updateMinDist(&minDistMatrix, resultMatrix, -1) // -1 so that length 1 paths are added

	// each iteration resultMatrix is changed to represent paths of length 1 greater than before
	// once all the nodes are connected, the loop will execute one more time
	// the value in insertions will be zero
	insertions := uint(math.MaxUint32)
	iterations := 0
	for insertions > 0 {
		timeStart := time.Now()
		resultMatrix = adjmat.Multiply(resultMatrix, adjacencyMatrix)
		insertions = updateMinDist(&minDistMatrix, resultMatrix, iterations)
		iterations++
		fmt.Printf("Completed loop %d (%v).\n", iterations, time.Now().Sub(timeStart))
	}

	return minDistMatrix
}

// updateMindist Edits minDistMatrix in place and places the minimum of the (i, j)th element of both matrices in (i, j)
// returns the number of elements that were changed (drawn from resultMatrix)
func updateMinDist(minDistMatrix *[][]uint16, resultMatrix *adjmat.AdjacencyMatrix, iterations int) uint {
	var insertions uint = 0
	rows, _ := resultMatrix.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < rows; j++ {
			if (*minDistMatrix)[i][j] == 0 && (*resultMatrix).At(i, j) != 0 {
				(*minDistMatrix)[i][j] = uint16(iterations + 2)
				insertions++
			}
		}
	}
	return insertions
}

func savePlot(fileName string, distanceDistribution map[uint16]uint16) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = fileName
	p.X.Label.Text = "Distance"
	p.Y.Label.Text = "Frequency"
	points := makePoints(distanceDistribution)
	err = plotutil.AddLinePoints(p, "Distribution", points)
	if err != nil {
		panic(err)
	}
	if err = p.Save(8*vg.Inch, 8*vg.Inch, "distance_distribution_"+fileName+".png"); err != nil {
		panic(err)
	}
}

func makePoints(distanceDistribution map[uint16]uint16) plotter.XYs {
	sortedDistances := makeSortedListOfDistances(distanceDistribution)
	points := make(plotter.XYs, len(distanceDistribution))
	// point := 0
	// for distance, frequency := range distanceDistribution {
	// 	points[point].X = float64(distance)
	// 	points[point].Y = float64(frequency)
	// 	point++
	// }
	for i := 0; i < len(sortedDistances); i++ {
		points[i].X = float64(sortedDistances[i])
		points[i].Y = float64(distanceDistribution[sortedDistances[i]])
	}
	return points
}

func makeSortedListOfDistances(distanceDistribution map[uint16]uint16) []uint16 {
	sortedDistances := make([]uint16, len(distanceDistribution))
	i := 0
	for distance := range distanceDistribution {
		sortedDistances[i] = distance
		i++
	}
	sort.Slice(sortedDistances, func(i, j int) bool { return sortedDistances[i] < sortedDistances[j] })
	return sortedDistances
}

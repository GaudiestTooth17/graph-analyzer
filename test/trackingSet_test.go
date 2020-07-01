package test

import (
	"testing"

	"gitlab.com/hcmi/graph-analyzer/trackingset"

	"gonum.org/v1/gonum/mat"
)

func TestAdd(t *testing.T) {
	rows, cols := 3, 3
	allZeros := mat.NewDense(rows, cols, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0})
	allOnes := mat.NewDense(rows, cols, []float64{1, 1, 1, 1, 1, 1, 1, 1, 1})
	diagonalEntries := mat.NewDense(rows, cols, []float64{1, 0, 0, 0, 1, 0, 0, 0, 1})
	reverseDiagonal := mat.NewDense(rows, cols, []float64{0, 0, 1, 0, 1, 0, 1, 0, 0})
	topRow := mat.NewDense(rows, cols, []float64{1, 1, 1, 0, 0, 0, 0, 0, 0})
	bottomRow := mat.NewDense(rows, cols, []float64{0, 0, 0, 0, 0, 0, 1, 1, 1})
	leftCol := mat.NewDense(rows, cols, []float64{1, 0, 0, 1, 0, 0, 1, 0, 0})
	middleCol := mat.NewDense(rows, cols, []float64{0, 1, 0, 0, 1, 0, 0, 1, 0})

	tSet := populateSet(allZeros)
	if int(tSet.Size()) != 0 {
		t.Error("Set was not empty.")
	}

	tSet = populateSet(allOnes)
	if int(tSet.Size()) != rows*cols {
		t.Errorf("Set did not contain %d elements.", rows*cols)
	}

	tSet = populateSet(diagonalEntries)
	if int(tSet.Size()) != rows {
		t.Errorf("Set did not contain %d entries.", rows)
	}
	tSet.Add(0, 0)
	if int(tSet.Size()) != rows {
		t.Errorf("Expected size: %d. Actual size: %d.", rows, int(tSet.Size()))
	}
	tSet.Add(0, 1)
	if int(tSet.Size()) != rows+1 {
		t.Errorf("Expected size: %d. Actual size: %d.", rows+1, int(tSet.Size()))
	}
	tSet.Add(1, 0)
	if int(tSet.Size()) != rows+2 {
		t.Errorf("Expected size: %d. Actual size: %d.", rows+2, int(tSet.Size()))
	}

	tSet = populateSet(diagonalEntries)
	tSet.AddAllNonzero(reverseDiagonal)
	if int(tSet.Size()) != 2*rows-1 {
		t.Errorf("Expected size: %d. Actual size: %d.", 2*rows-1, int(tSet.Size()))
	}
	tSet.AddAllNonzero(reverseDiagonal)
	if int(tSet.Size()) != 2*rows-1 {
		t.Errorf("Size changed. Expected size: %d. Actual size: %d.", 2*rows-1, int(tSet.Size()))
	}
	tSet.AddAllNonzero(allOnes)
	if int(tSet.Size()) != rows*cols {
		t.Errorf("Expected size: %d. Actual size: %d.", rows*cols, int(tSet.Size()))
	}

	tSet = populateSet(topRow)
	tSet.AddAllNonzero(bottomRow)
	if int(tSet.Size()) != 2*rows {
		t.Errorf("Expected size: %d. Actual size: %d.", 2*rows, int(tSet.Size()))
	}

	tSet = populateSet(middleCol)
	tSet.AddAllNonzero(leftCol)
	if int(tSet.Size()) != 2*cols {
		t.Errorf("Expected size: %d. Actual size: %d.", 2*cols, int(tSet.Size()))
	}
}

func populateSet(matrix mat.Matrix) trackingset.Set {
	rows, cols := matrix.Dims()
	tSet := trackingset.NewSet(rows, cols)
	tSet.AddAllNonzero(matrix)
	return tSet
}

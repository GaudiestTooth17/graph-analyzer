package trackingset

import (
	"math"

	"gitlab.com/hcmi/graph-analyzer/adjmat"
)

type void struct{}

// Set ... A set to keep track of what entries in a matrix are nonzero
type Set struct {
	data         []uint32
	rows         uint32
	cols         uint32
	elementCount uint32
}

// AddAllNonzero ... Add all the nonzero entries of a matrix to the set
func (s *Set) AddAllNonzero(matrix *adjmat.AdjacencyMatrix) {
	rows, cols := matrix.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if matrix.At(i, j) != 0 {
				(*s).Add(i, j)
			}
		}
	}
}

// NewSet ... Make a new TrackingSet of specified size
func NewSet(rows, cols int) Set {
	lenData := (rows * cols) / 32
	if (rows*cols)%32 != 0 {
		lenData++
	}
	return Set{rows: uint32(rows), cols: uint32(cols), data: make([]uint32, lenData), elementCount: 0}
}

// Add ... Add the given entry to the set
func (s *Set) Add(row, col int) {
	neighborhood := ((*s).rows*uint32(row) + uint32(col)) / 32
	house := ((*s).rows*uint32(row) + uint32(col)) % 32
	newResident := uint32(1 << house)
	currentResidents := (*s).data[neighborhood]
	// create a mask to check if the data already contains this entry
	if (newResident^math.MaxUint32)&currentResidents == currentResidents {
		(*s).data[neighborhood] += newResident
		(*s).elementCount++
	}
}

// Contains ... true if the set contains the given coordinates
func (s *Set) Contains(row, col int) bool {
	neighborhood := ((*s).rows*uint32(row) + uint32(col)) / 32
	house := ((*s).rows*uint32(row) + uint32(col)) % 32
	query := uint32(1 << house)
	currentResidents := (*s).data[neighborhood]
	// create a mask to check if the data already contains this entry
	if (query^math.MaxUint32)&currentResidents != currentResidents {
		return false
	}
	return true
}

// Size ... return the number of elements in the set
func (s *Set) Size() uint32 {
	return (*s).elementCount
}

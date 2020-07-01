package adjmat

import (
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas32"
)

// AdjacencyMatrix stores representation of a weighted undirected graph
type AdjacencyMatrix struct {
	mat blas32.Symmetric
}

// Multiply returns A * B
func Multiply(A, B *AdjacencyMatrix) *AdjacencyMatrix {
	C := New(A.mat.N, make([]float32, A.mat.N*A.mat.N))
	cMat := blas32.General{Rows: C.mat.N, Cols: C.mat.N, Stride: C.mat.Stride, Data: C.mat.Data}
	aMat := blas32.General{Rows: A.mat.N, Cols: A.mat.N, Stride: A.mat.Stride, Data: A.mat.Data}
	blas32.Symm(blas.Right, 1, B.mat, aMat, 0, cMat)
	return C
}

// New returns a pointer to an n x n AdjacencyMatrix
func New(n int, data []float32) *AdjacencyMatrix {
	return &AdjacencyMatrix{
		mat: blas32.Symmetric{
			N:      n,
			Stride: n,
			Data:   data,
			Uplo:   blas.Upper,
		},
	}
}

// CopyOf returns a copy of A
func CopyOf(A *AdjacencyMatrix) *AdjacencyMatrix {
	return New(A.mat.N, A.mat.Data)
}

//At returns the element at row i and column j
func (a *AdjacencyMatrix) At(i, j int) float32 {
	if uint(i) >= uint(a.mat.N) {
		panic("Row out of bounds!")
	}
	if uint(j) >= uint(a.mat.N) {
		panic("Column out of bounds!")
	}
	if i > j {
		i, j = j, i
	}
	return a.mat.Data[i*a.mat.Stride+j]
}

// Dims returns the number of rows and columns in the matrix.
func (a *AdjacencyMatrix) Dims() (r, c int) {
	return a.mat.N, a.mat.N
}

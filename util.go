package main

import "fmt"

func printMatrix(matrix [][]uint16) {
	fmt.Println("[")
	for row := 0; row < len(matrix)-1; row++ {
		fmt.Println(matrix[row])
	}
	fmt.Printf("%v\n]\n", matrix[len(matrix)-1])
}

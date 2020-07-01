package main

// createGraphFromComponent takes in a "set" of nodes to keep in graph and returns a new matrix
// with just the nodes in the set. Their indices are, of course, changed.
func createGraphFromComponent(component map[int]void, graph [][]uint16) [][]uint16 {
	// create graph
	newGraph := make([][]uint16, len(component))
	for i := 0; i < len(component); i++ {
		newGraph[i] = make([]uint16, len(component))
	}

	// populate graph
	oldToNew := makeOldToNewIndexMap(component)
	for node := range component {
		for neighbor := range component {
			newGraph[oldToNew[node]][oldToNew[neighbor]] = graph[node][neighbor]
		}
	}

	return newGraph
}

func makeOldToNewIndexMap(component map[int]void) map[int]int {
	oldToNew := make(map[int]int, len(component))
	currentIndex := 0
	for node := range component {
		oldToNew[node] = currentIndex
		currentIndex++
	}
	return oldToNew
}

package main

func findComponents(adjMatrix [][]uint16) []map[int]void {
	components := make([]map[int]void, 0)
	visited := make(map[int]bool)

	for node := 0; node < len(adjMatrix[0]); node++ {
		if visited[node] {
			continue
		}

		component := explore(node, adjMatrix, visited)
		components = append(components, component)
		//union visited and component
		for n := range component {
			visited[n] = true
		}
	}

	return components
}

func explore(node int, adjMatrix [][]uint16, visited map[int]bool) map[int]void {
	var exists void
	component := make(map[int]void, 0)
	component[node] = exists
	visited[node] = true

	for neighbor := 0; neighbor < len(adjMatrix[node]); neighbor++ {
		if adjMatrix[node][neighbor] != 0 && !visited[neighbor] {
			newNodes := explore(neighbor, adjMatrix, visited)
			//union
			for newNode := range newNodes {
				component[newNode] = exists
			}
		}
	}

	return component
}

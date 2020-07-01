import sys
import numpy as np
from typing import Set, List, Tuple
from scipy.sparse.csgraph import floyd_warshall
import matplotlib.pyplot as plt
import time


def filter_adjacency_matrix(adj_matrix: np.ndarray, nodes_to_keep: Set[int])\
                            -> np.ndarray:
    new_adj_matrix = np.zeros((len(nodes_to_keep), len(nodes_to_keep)),
                              dtype=np.uint16)
    for i, node in enumerate(nodes_to_keep):
        for j, neighbor in enumerate(nodes_to_keep):
            new_adj_matrix[i][j] = adj_matrix[node][neighbor]
    return new_adj_matrix


def calculate_degree_distribution(distance_matrix: np.ndarray)\
                                  -> Tuple[List[int], List[int]]:
    degree_to_frequency = dict()
    for row in distance_matrix:
        for degree in row:
            if degree not in degree_to_frequency:
                degree_to_frequency[degree] = 0
            degree_to_frequency[degree] += 1

    return [degree for degree in degree_to_frequency.keys()],\
           [frequency for frequency in degree_to_frequency.values()]


def main():
    if len(sys.argv) < 3:
        print('Usage:', sys.argv[0], '[graph file] [file of nodes to include]')
        return

    graph_file_name = sys.argv[1]
    node_file_name = sys.argv[2]

    time_start = time.time()
    adj_matrix = np.loadtxt(graph_file_name, dtype=np.uint16)
    adj_matrix = np.where(adj_matrix > 0, 1, 0)
    print('Loaded matrix in', time.time() - time_start)

    with open(node_file_name, 'r') as file:
        nodes_to_keep = {int(line) for line in file}

    time_start = time.time()
    adj_matrix = filter_adjacency_matrix(adj_matrix, nodes_to_keep)
    print('filtered matrix in', time.time() - time_start)
    time_start = time.time()
    distance_matrix = floyd_warshall(csgraph=adj_matrix, directed=False,
                                     return_predecessors=False)
    print('Completed floyd warshall in', time.time() - time_start)
    np.save('distance_matrix.npy', distance_matrix)
    # degrees, frequencies = calculate_degree_distribution(distance_matrix)
    # plt.plot(degrees, frequencies)

    # remove all nodes not in nodes to keep from adj_matrix
    # see go code for algorithm
    # do floyd's
    # make distribution


if __name__ == "__main__":
    main()

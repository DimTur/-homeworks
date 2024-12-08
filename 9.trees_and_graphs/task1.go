// Реализуйте BFS алгоритм в представлении матрицы стоимости

package main

import (
	"fmt"
	"math"
)

type MtxList struct {
	adjMatrix [][]int
}

type Dist struct {
	Dist     int
	EdgeFrom int
}

var mtx = [][]int{
	{0, 2, 3, 0, 0, 0},
	{2, 0, 0, 1, 0, 0},
	{3, 0, 0, 0, 0, 1},
	{0, 1, 0, 0, 0, 2},
	{0, 1, 0, 0, 0, 0},
	{0, 0, 1, 2, 0, 0},
}

func (m *MtxList) BFSShort(start, target int) int {
	d := make([]Dist, len(m.adjMatrix))
	for i := range d {
		d[i].Dist = math.MaxInt32
		d[i].EdgeFrom = -1
	}
	d[start].Dist = 0

	var queue []int
	queue = append(queue, start)

	for len(queue) != 0 {
		currentVert := queue[0]

		for adjVert, weight := range m.adjMatrix[currentVert] {
			if weight > 0 && d[currentVert].Dist+weight < d[adjVert].Dist {
				d[adjVert].Dist = d[currentVert].Dist + weight
				queue = append(queue, adjVert)
			}
		}
		queue = queue[1:]
	}

	if d[target].EdgeFrom == math.MaxInt32 {
		return -1
	}
	return d[target].Dist
}

func main() {
	start := 4
	target := 0

	graph := MtxList{adjMatrix: mtx}

	fmt.Println(graph.BFSShort(start, target))
}

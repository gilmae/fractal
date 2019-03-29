package fractal

import (
	"sync"
)

var counter = struct {
	sync.RWMutex
	visited map[Key]bool
}{visited: make(map[Key]bool)}

func Find_Escapee(points_map map[Key]CalculatedPoint) Key {
	for k, v := range points_map {
		if v.Escaped {
			return k
		}
	}

	return Key{-1, -1}
}

func check_is_edge(k Key, edgePoints chan CalculatedPoint, points_map map[Key]CalculatedPoint, width int, height int) {
	if k.X < 0 || k.X >= width || k.Y < 0 || k.Y > height {
		return
	}

	if _, ok := counter.visited[k]; ok {
		return
	}
	counter.Lock()
	counter.visited[k] = true
	counter.Unlock()

	if points_map[k].Escaped {
		edgePoints <- points_map[k]
		return
	}

	for ii := -1; ii < 2; ii++ {
		for jj := -1; jj < 2; jj++ {
			if ii != 0 || jj != 00 {
				check_is_edge(Key{k.X + ii, k.Y + jj}, edgePoints, points_map, width, height)
			}
		}
	}
}

func Find_Edges(edgePoints chan CalculatedPoint, points_map map[Key]CalculatedPoint, width int, height int) {
	// scan for a non-escaped pixel
	var p = Find_Escapee(points_map)

	if p.X >= 0 {
		check_is_edge(p, edgePoints, points_map, width, height)
	}
}

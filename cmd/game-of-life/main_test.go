package main

import "testing"

func TestSetFixedNeighborhood(t *testing.T) {
	// (w *World) setFixedNeighborhood(x, y int, fixed []bool) map[int]struct{}
	width := 3
	height := 3
	w := World{
		area:    make([]bool, width*height),
		mem:     make([][]int, width*height),
		width:   width,
		height:  height,
		paused:  true,
		forward: false,
	}
	fixed := make([]bool, width*height)
	fixed[2] = true
	changes := w.setFixedNeighborhood(1, 1, fixed)
	if _, ok := changes[2]; ok {
		t.Errorf("should not have changed 2")
	}
}

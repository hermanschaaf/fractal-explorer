package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	screenWidth  = 9
	screenHeight = 9
	SearchDepth  = 1
)

// World represents the game state.
type World struct {
	area    []bool
	width   int
	height  int
	paused  bool
	forward bool
}

// NewWorld creates a new world.
func NewWorld(width, height int, maxInitLiveCells int) *World {
	w := &World{
		area:    make([]bool, width*height),
		width:   width,
		height:  height,
		paused:  true,
		forward: false,
	}
	//w.init(maxInitLiveCells)
	return w
}

// init inits world with a random state.
func (w *World) init(maxLiveCells int) {
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.area[y*w.width+x] = true
	}
}

// Update game state by one tick.
func (w *World) Update() {
	if w.paused {
		return
	}

	if w.forward {
		w.gameOfLife()
	} else {
		w.reverseGameOfLife()
	}
}

// 				// rule 1. Any live cell with fewer than two live neighbours
//				// dies, as if caused by under-population.
//				// rule 2. Any live cell with two or three live neighbours
//				// lives on to the next generation.
//				// rule 3. Any live cell with more than three live neighbours
//				// dies, as if by over-population.
//				// rule 4. Any dead cell with exactly three live neighbours
//				// becomes a live cell, as if by reproduction.
func (w *World) gameOfLife() {
	width := w.width
	height := w.height
	next := make([]bool, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pop := neighbourCount(w.area, width, height, x, y)
			switch {
			case pop < 2:
				next[y*width+x] = false

			case (pop == 2 || pop == 3) && w.area[y*width+x]:
				next[y*width+x] = true

			case pop > 3:
				next[y*width+x] = false

			case pop == 3:
				next[y*width+x] = true
			}
		}
	}
	w.area = next
}

func (w *World) reverseGameOfLife() {
	width := w.width
	height := w.height
	prev := make([]bool, width*height)
	//solved := w.solve(prev, w.area, 0, 0)
	log.Println("solving...")
	solved := w.solve2(prev, w.area, SearchDepth)
	log.Println("solved:", solved)
	if solved {
		w.area = prev
	}
}

// figure out valid `prev` that is possible given `next` using backtracking
func (w *World) solve2(prev, next []bool, depth int) bool {
	if depth <= 0 {
		return true
	}
	log.Printf("solving for depth %d...", depth)

	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			l := y*w.width + x
			if next[l] {
				fmt.Printf("X")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")

	// build list of live cells in next iteration
	live := make([]int, 0)
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			i := y*w.width + x
			if next[i] {
				live = append(live, i)
			}
		}
	}
	fixed := make([]bool, w.width*w.height)
	return w.solve2inner(prev, next, fixed, live, 0, 0, depth)
}

func (w *World) solve2inner(prev, next, fixed []bool, live []int, i int, fixedUpTo int, depth int) bool {
	if i >= len(live) {
		return w.allok(prev, next)
	}
	log.Println(i, fixedUpTo)
	//
	//for y := 0; y < w.height; y++ {
	//	for x := 0; x < w.width; x++ {
	//		l := y*w.width + x
	//		if live[i] == l {
	//			fmt.Printf("i")
	//		} else if prev[l] {
	//			fmt.Printf("X")
	//		} else if fixed[l] {
	//			fmt.Printf("o")
	//		} else {
	//			fmt.Printf(".")
	//		}
	//	}
	//	fmt.Println("")
	//}
	//fmt.Println("")

	// live cells must be one of the following in the previous step:
	//  1. live with 2 neighbors
	//  2. live with 3 neighbors
	//  3. dead with 3 neighbors
	l := live[i]
	x, y := l%w.width, l/w.width
	pop := neighbourCount(prev, w.width, w.height, x, y)
	if pop >= 4 {
		return false
	}

	// check fixed area for violations
	if !w.fixedok(prev, fixed, next, l) {
		return false
	}

	changed := -1
	// if alive in next step, try and make this cell live in previous
	if next[l] && !fixed[l] && !prev[l] {
		prev[l] = true
		changed = l
	}

	if prev[l] {
		// if cell is alive
		if pop <= 1 {
			// try and add a neighbor for case 1
			for z := fixedUpTo; z < 8; z++ {
				n := w.addNeighbor(x, y, prev, fixed, z)
				if n == -1 {
					w.undoChange(changed, prev)
					return false
				}

				solved := w.solve2inner(prev, next, fixed, live, i, z+1, depth)
				if solved {
					nextDepth := make([]bool, len(prev))
					solvedNextDepth := w.solve2(nextDepth, prev, depth-1)
					if solvedNextDepth {
						return true
					}
				}
				// undo adding of neighbor
				w.undoAddNeighbor(prev, n)
			}
		} else if pop == 2 {
			// try and solve with no further changes (except for fixed cells)
			changes := w.setFixedNeighborhood(x, y, fixed)
			solved := w.solve2inner(prev, next, fixed, live, i+1, 0, depth)
			if solved {
				nextDepth := make([]bool, len(prev))
				solvedNextDepth := w.solve2(nextDepth, prev, depth-1)
				if solvedNextDepth {
					return true
				}
			}
			w.undoFixedNeighborhood(fixed, changes)

			// also try adding one more neighbor for case 2
			for z := fixedUpTo; z < 8; z++ {
				n := w.addNeighbor(x, y, prev, fixed, z)
				if n == -1 {
					w.undoChange(changed, prev)
					return false
				}
				solved = w.solve2inner(prev, next, fixed, live, i, z+1, depth)
				if solved {
					nextDepth := make([]bool, len(prev))
					solvedNextDepth := w.solve2(nextDepth, prev, depth-1)
					if solvedNextDepth {
						return true
					}
				}
				// undo adding of neighbor
				w.undoAddNeighbor(prev, n)
			}
		} else if pop == 3 {
			// only remaining case since we checked pop >= 4 before
			// try and solve with no further changes (except for fixed cells)
			changes := w.setFixedNeighborhood(x, y, fixed)
			solved := w.solve2inner(prev, next, fixed, live, i+1, 0, depth)
			if solved {
				nextDepth := make([]bool, len(prev))
				solvedNextDepth := w.solve2(nextDepth, prev, depth-1)
				if solvedNextDepth {
					return true
				}
			}
			w.undoFixedNeighborhood(fixed, changes)
		}
	}

	w.undoChange(changed, prev)

	// if alive in next step, try and make this cell dead in previous (case 3)
	if next[l] && !fixed[l] && prev[l] {
		prev[l] = false
		changed = l
	}

	if !prev[l] {
		// if cell is dead
		if pop <= 2 {
			// try and add a neighbor
			for z := fixedUpTo; z < 8; z++ {
				n := w.addNeighbor(x, y, prev, fixed, z)
				if n == -1 {
					w.undoChange(changed, prev)
					return false
				}
				solved := w.solve2inner(prev, next, fixed, live, i, z+1, depth)
				if solved {
					nextDepth := make([]bool, len(prev))
					solvedNextDepth := w.solve2(nextDepth, prev, depth-1)
					if solvedNextDepth {
						return true
					}
				}
				// undo adding of neighbor
				w.undoAddNeighbor(prev, n)
			}
		} else if pop == 3 {
			// only remaining case since we checked pop >= 4 before
			// try and solve with no further changes (except for fixed cells)
			changes := w.setFixedNeighborhood(x, y, fixed)
			solved := w.solve2inner(prev, next, fixed, live, i+1, 0, depth)
			if solved {
				nextDepth := make([]bool, len(prev))
				solvedNextDepth := w.solve2(nextDepth, prev, depth-1)
				if solvedNextDepth {
					return true
				}
			}
			w.undoFixedNeighborhood(fixed, changes)
		}
	}

	// ok, we didn't find a solution, so we give up at this point
	w.undoChange(changed, prev)

	return false
}

func (w *World) undoChange(l int, prev []bool) {
	if l == -1 {
		return
	}
	prev[l] = !prev[l]
}

func (w *World) fixNeighborhoodUpTo(x, y, i int, fixed []bool) map[int]struct{} {
	changes := make(map[int]struct{})
	ind := 0
	for xi := x - 1; xi <= x+1; xi++ {
		for yi := y - 1; yi <= y+1; yi++ {
			if (xi == x && yi == y) || xi < 0 || yi < 0 || xi >= w.width || yi >= w.height {
				continue
			}
			l := yi*w.width + xi
			if ind < i && !fixed[l] {
				fixed[l] = true
				changes[l] = struct{}{}
			}
			ind++
		}
	}
	return changes
}

func (w *World) setFixedNeighborhood(x, y int, fixed []bool) map[int]struct{} {
	return w.fixNeighborhoodUpTo(x, y, 8, fixed)
}

func (w *World) undoFixedNeighborhood(fixed []bool, changes map[int]struct{}) {
	for c := range changes {
		fixed[c] = false
	}
}

func (w *World) addNeighbor(x, y int, ar []bool, fixed []bool, fixedUpTo int) int {
	ind := 0
	for xi := x - 1; xi <= x+1; xi++ {
		for yi := y - 1; yi <= y+1; yi++ {
			if (xi == x && yi == y) || xi < 0 || yi < 0 || xi >= w.width || yi >= w.height {
				continue
			}
			if ind < fixedUpTo {
				ind++
				continue
			}

			l := yi*w.width + xi
			if fixed[l] {
				continue
			}
			if ar[l] {
				continue
			}
			ar[l] = true
			return l
		}
	}
	return -1
}

func (w *World) undoAddNeighbor(ar []bool, l int) {
	ar[l] = false
}

// figure out valid `prev` that is possible given `next` using backtracking
func (w *World) solve(prev, next []bool, r, c int) bool {
	if r >= w.height {
		return w.allok(prev, next)
	}

	prev[r*w.width+c] = false

	// check if constraints at (r-1, c-1) are being violated
	pr, pc := r-1, c-1
	ok := w.ok(prev, next, pr, pc)

	nr, nc := r, c+1
	if nc >= w.width {
		nr, nc = r+1, 0
	}

	if ok {
		solved := w.solve(prev, next, nr, nc)
		if solved {
			return solved
		}
	}

	// set value at index to 1 and check if it violates any constraints
	prev[r*w.width+c] = true
	ok = w.ok(prev, next, pr, pc)

	if ok {
		solved := w.solve(prev, next, nr, nc)
		if solved {
			return solved
		}
	}

	// undo our changes
	prev[r*w.width+c] = false
	return false
}

// fixedok checks if any violations occur in areas that are marked as "fixed". Anything to the
// upper left of currentIndex and before will also be taken to be fixed.
func (w *World) fixedok(prev, fixed, wantNext []bool, currentIndex int) bool {
	width := w.width
	height := w.height

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			l := y*width + x
			isFixed := (l < currentIndex-width-1) || (fixed[l] && neighbourCount(fixed, width, height, x, y) == 8)
			// this cell is completely surrounded by fixed cells
			if isFixed {
				pop := neighbourCount(prev, width, height, x, y)
				if wantNext[l] {
					// alive in next; should have (pop == 2 and alive) or (pop == 3) in prev
					if (prev[l] && pop == 2) || (pop == 3) {
						continue
					}
					return false
				} else {
					// dead in next; should have pop == 0, 1 or >= 4, or be dead with pop == 2
					if (!prev[l] && pop == 2) || (pop <= 1) || (pop >= 4) {
						continue
					}
					return false
				}
			}
		}
	}
	return true
}

func (w *World) allok(prev, wantNext []bool) bool {
	width := w.width
	height := w.height
	next := make([]bool, width*height)

	allSame := true
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pop := neighbourCount(prev, width, height, x, y)
			switch {
			case pop < 2:
				next[y*width+x] = false
			case (pop == 2 || pop == 3) && prev[y*width+x]:
				next[y*width+x] = true
			case pop > 3:
				next[y*width+x] = false
			case pop == 3:
				next[y*width+x] = true
			}
			if wantNext[y*width+x] != next[y*width+x] {
				return false
			}
			if prev[y*width+x] != next[y*width+x] {
				allSame = false
			}
		}
	}
	return !allSame
}

func (w *World) ok(prev, next []bool, r, c int) bool {
	if r >= 0 && c >= 0 && r < w.height && c < w.width {
		pop := neighbourCount(prev, w.width, w.height, r, c)
		n := next[w.width*r+c]
		if n && (pop < 2 || pop > 3) {
			// impossible: cell is alive in next iteration but had too few neighbors to support it
			return false
		} else if !n && (pop == 3) {
			// impossible: cell is dead in next iteration but it had exactly 3 neighbors
			return false
		}
	}
	return true
}

func sum(v []int) int {
	s := 0
	for i := range v {
		s += v[i]
	}
	return s
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	for i, v := range w.area {
		if v {
			pix[4*i] = 0xff
			pix[4*i+1] = 0xff
			pix[4*i+2] = 0xff
			pix[4*i+3] = 0xff
		} else {
			pix[4*i] = 0
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0
		}
	}
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// neighbourCount calculates the Moore neighborhood of (x, y).
func neighbourCount(a []bool, width, height, x, y int) int {
	c := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}
			x2 := x + i
			y2 := y + j
			if x2 < 0 || y2 < 0 || width <= x2 || height <= y2 {
				continue
			}
			if a[y2*width+x2] {
				c++
			}
		}
	}
	return c
}

type Game struct {
	world  *World
	pixels []byte
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()

	if g.world.paused && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x := max(0, min(x, screenWidth))
		y := max(0, min(y, screenHeight))
		g.world.area[y*screenWidth+x] = true
	}
	if g.world.paused && ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x := max(0, min(x, screenWidth))
		y := max(0, min(y, screenHeight))
		g.world.area[y*screenWidth+x] = false
	}

	if repeatingKeyPressed(ebiten.KeyP) {
		g.world.paused = !g.world.paused
	}
	if repeatingKeyPressed(ebiten.KeyF) {
		g.world.forward = !g.world.forward
	}

	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.world.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)
	log.Printf("Paused: %v, Forward: %v", g.world.paused, g.world.forward)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Paused: %v", g.world.paused))
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	//if d >= delay && (d-delay)%interval == 0 {
	//	return true
	//}
	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		world: NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
	}
	ebiten.SetMaxTPS(10)
	ebiten.SetWindowSize(screenWidth*20, screenHeight*20)
	ebiten.SetWindowTitle("Modified Game of Life")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

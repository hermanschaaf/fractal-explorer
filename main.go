package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"fractal-explorer/linalg"
)

// World represents the state.
type World struct {
	width  int
	height int
	center linalg.Vec3
	zoom   float64
}

// NewWorld creates a new world.
func NewWorld(width, height int) *World {
	w := &World{
		width:  width,
		height: height,
		center: linalg.Vec3{},
		zoom:   12.0,
	}
	return w
}

// Update game state by one tick.
func (w *World) Update() {
	w.center = w.center.Add(linalg.Vec3{Z: 0.001})
	w.zoom *= 0.999
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	start := w.center.Add(linalg.Vec3{X: -w.zoom / 2, Y: -w.zoom / 2, Z: 0})
	step := linalg.Vec3{
		X: w.zoom / float64(w.width),
		Y: w.zoom / float64(w.height),
		Z: 0.0,
	}
	p := start
	maxIter := 20
	a, b, c := -0.5, 1., -0.1
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			ans := iterate(a, b, c, p, maxIter)
			if ans == 0 {
				w.setPixel(pix, x, y, color.Black)
			} else {
				w.setPixel(pix, x, y, color.RGBA{R: uint8(255 * (float64(ans) / float64(maxIter))), G: 0, B: 0, A: 0xff})
			}
			p.X += step.X
		}
		p.X = start.X
		p.Y += step.Y
	}
}

func (w *World) setPixel(pix []byte, x, y int, c color.Color) {
	i := y*w.width + x
	r, g, b, a := c.RGBA()
	pix[4*i] = byte(r)
	pix[4*i+1] = byte(g)
	pix[4*i+2] = byte(b)
	pix[4*i+3] = byte(a)
}

// f_0(x, y) = a x^2 + b xy + c y^2
// f_n+1(x, y) = f_n(x, y)^2 + (a x^2 + b xy + c y^2)
func iterate(a, b, c float64, position linalg.Vec3, maxIter int) int {
	q := a*position.X*position.X + b*position.X*position.Y + c*position.Y*position.Y
	z := 0.0
	for i := 0; i < maxIter; i++ {
		z = z*z + q
		if z > 100000000 {
			return i + 1
		}
	}
	return 0
}

const (
	screenWidth  = 320
	screenHeight = 240
)

type Game struct {
	world  *World
	pixels []byte
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.world.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		world: NewWorld(screenWidth, screenHeight),
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Herman's Fractal Explorer")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

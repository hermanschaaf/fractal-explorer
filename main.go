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
		zoom:   4.0,
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
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			z := linalg.Vec3{X: 0, Y: 0, Z: 0}
			ans := iterate(z, p, maxIter)
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

func iterate(z, position linalg.Vec3, maxIter int) int {
	for i := 0; i < maxIter; i++ {
		z = z.Multiply2D(z).Add(position)
		l := z.Length()
		if l > 2 {
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

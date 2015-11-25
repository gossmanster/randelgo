// Package randelbrot allows a user to create images of the famous Mandelbrot set fractal. More specifically,
// it has the ability to randomly navigate the Mandelbrot Set, generating interesting images
package randelbrot

import (
	"math/rand"
)

// RandelbrotServer holds the state for a server that can explore the Mandelbrot Set
type RandelbrotServer struct {
	random *rand.Rand
	latest *MandelbrotSet
}

// RenderToBuffer actually draws a set into a PixelBuffer 
func RenderToBuffer(buffer *PixelBuffer, set *MandelbrotSet) {
	renderer := new(Renderer)
	maxCount := set.EstimateMaxCount()

	bandMap := NewLogarithmicBandMap(maxCount, 32.0)

	renderer.RenderByCrawling(buffer, set, bandMap, maxCount)
}

func (server *RandelbrotServer) randomChild(set *MandelbrotSet) *MandelbrotSet {
	newSide := (server.random.Float64() * set.Side / 4.5) + set.Side/6
	newCX := ((server.random.Float64() - 0.5) * set.Side / 2) + set.CenterX
	newCY := ((server.random.Float64() - 0.5) * set.Side / 2) + set.CenterY
	newSet := new(MandelbrotSet)
	newSet.CenterY = newCX
	newSet.Side = newSide
	newSet.CenterY = newCY

	return newSet
}

// NewRandelbrotServer creates a server
func NewRandelbrotServer(random *rand.Rand) (server *RandelbrotServer) {
	server = new(RandelbrotServer)
	server.random = random
	server.latest = new(MandelbrotSet)
	server.latest.CenterY = 0.0
	server.latest.CenterX = -0.5
	server.latest.Side = 2.5

	return server
}

// RenderNext creates the next image in a sequence
func (server *RandelbrotServer) RenderNext(buffer *PixelBuffer) {
	RenderToBuffer(buffer, server.latest)
	server.latest = server.randomChild(server.latest)
}

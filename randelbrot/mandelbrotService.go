// Package randelbrot allows a user to create images of the famous Mandelbrot set fractal. More specifically,
// it has the ability to randomly navigate the Mandelbrot Set, generating interesting images
package randelbrot

import (
	log "github.com/Sirupsen/logrus"
	"math"
	"math/rand"
)

// RandelbrotServer holds the state for a server that can explore the Mandelbrot Set
type RandelbrotServer struct {
	random *rand.Rand
	latest *MandelbrotSet
	futures *priorityQueue
}

// RenderToBuffer actually draws a set into a PixelBuffer
func RenderToBuffer(buffer *PixelBuffer, set *MandelbrotSet) {
	renderer := new(Renderer)
	maxCount := set.EstimateMaxCount()

	bandMap := NewLogarithmicBandMap(maxCount, 32.0)

	renderer.RenderByCrawling(buffer, set, bandMap, maxCount)
}

// NewRandelbrotServer creates a server
func NewRandelbrotServer(random *rand.Rand) (server *RandelbrotServer) {
	server = new(RandelbrotServer)
	server.random = random
	server.latest = new(MandelbrotSet)
	server.latest.CenterY = 0.0
	server.latest.CenterX = -0.5
	server.latest.Side = 2.5
	
	server.futures = newPriorityQueue(1000)

	return server
}

// RenderNext creates the next image in a sequence
func (server *RandelbrotServer) RenderNext(buffer *PixelBuffer) {
	RenderToBuffer(buffer, server.latest)
	candidates := server.generateCandidates(server.latest)

	for i := 0; i < len(candidates); i++ {
		b := evaluateBeauty(candidates[i])
		server.futures.push(candidates[i], b)
	}
	server.latest = server.futures.pop()
	log.WithFields(log.Fields{
		"latestCX": server.latest.CenterX,
		"latestCY": server.latest.CenterY,
	}).Info("RenderNext Prepped")
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

func (server *RandelbrotServer) generateCandidates(set *MandelbrotSet) []*MandelbrotSet {
	count := 20
	retval := make([]*MandelbrotSet, count)
	for i := 0; i < count; i++ {
		retval[i] = server.randomChild(set)
	}

	return retval
}

func evaluateBeauty(set *MandelbrotSet) (evaluation float64) {
	bufferSize := 100
	renderer := new(Renderer)
	maxCount := set.EstimateMaxCount()

	bandMap := NewLogarithmicBandMap(maxCount, 30.0)
	buffer := NewPixelBuffer(bufferSize, bufferSize)

	bandCount := renderer.RenderByCrawling(buffer, set, bandMap, maxCount)
	evaluation = float64(bandCount)

	h := evaluateBuffer(buffer)

	pointsInSet := h.getValue(-1)
	if pointsInSet > 0 {
		// All black isn't pretty
		if pointsInSet == (bufferSize * bufferSize) {
			evaluation = math.MinInt64
		}
		// But some black is good
		evaluation *= 2

		// But not too much black
		r := float64(bufferSize*bufferSize) / float64(pointsInSet) / 200.0
		evaluation += r
		log.WithFields(log.Fields{
			"pointsInSet": pointsInSet,
			"blackRatio":  r,
			"evaluation":  evaluation,
			"maxCount":    maxCount,
		}).Info("adjusing beauty")
	}

	// More colors good
	evaluation += float64(h.valueCount()) / 3

	// Less max count good
	evaluation -= float64(maxCount) / 100

	log.WithFields(log.Fields{
		"colorValues": h.valueCount(),
		"evaluation":  evaluation,
	}).Info("Returning from evaluateBeauty")
	return evaluation
}

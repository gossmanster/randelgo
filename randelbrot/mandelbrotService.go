// Package randelbrot allows a user to create images of the famous Mandelbrot set fractal. More specifically,
// it has the ability to randomly navigate the Mandelbrot Set, generating interesting images
package randelbrot

import (
	log "github.com/Sirupsen/logrus"
	"math"
	"math/rand"
	"time"
)

// RandelbrotServer holds the state for a server that can explore the Mandelbrot Set
type RandelbrotServer struct {
	random  *rand.Rand
	root    *MandelbrotSet
	latest  *MandelbrotSet
	futures *priorityQueue
}

// RenderToBuffer actually draws a set into a PixelBuffer
func RenderToBuffer(buffer *PixelBuffer, set *MandelbrotSet) {
	log.WithFields(log.Fields{
		"CX":set.CenterX,
		"CY": set.CenterY,
		"Side": set.Side,
	}).Info("RenderToBuffer")
	renderer := new(Renderer)
	maxCount := set.EstimateMaxCount(150)

	bandMap := NewLogarithmicBandMap(maxCount, 50.0)

	renderer.RenderByCrawling(buffer, set, bandMap, maxCount)
}

// NewRandelbrotServer creates a server
func NewRandelbrotServer(random *rand.Rand) (server *RandelbrotServer) {
	server = new(RandelbrotServer)
	server.random = random
	server.root = new(MandelbrotSet)
	server.root.CenterY = 0.0
	server.root.CenterX = -0.5
	server.root.Side = 2.5
	server.latest = server.root

	server.futures = newPriorityQueue(1000)
	for i := 0; i < 200; i++ {
		candidate := server.randomChild(server.root)
		b := multiLevelEvaluation(candidate)
		server.futures.push(candidate, b)
	}

	return server
}

// RenderNext creates the next image in a sequence
func (server *RandelbrotServer) RenderNext(buffer *PixelBuffer) {
	log.Info("RenderNext")
	t := time.Now()
	RenderToBuffer(buffer, server.latest)
	log.Info("Evaluating Candidates")
	candidates := server.generateCandidates(server.latest)

	for i := 0; i < len(candidates); i++ {
		b := multiLevelEvaluation(candidates[i])
		server.futures.push(candidates[i], b)
	}
	server.latest = server.futures.pop()
	log.WithFields(log.Fields{
		"Duration": time.Since(t),
	}).Info("RenderNext done")
}

func (server *RandelbrotServer) randomChild(set *MandelbrotSet) *MandelbrotSet {
	newSide := (server.random.Float64() * set.Side / 4) + set.Side/6
	newCX := ((server.random.Float64() - 0.5) * set.Side / 1.2) + set.CenterX
	newCY := ((server.random.Float64() - 0.5) * set.Side / 1.2) + set.CenterY
	newSet := new(MandelbrotSet)
	newSet.CenterX = newCX
	newSet.Side = newSide
	newSet.CenterY = newCY
	return newSet
}

func (server *RandelbrotServer) generateCandidates(set *MandelbrotSet) []*MandelbrotSet {
	count := 6
	retval := make([]*MandelbrotSet, 0)
	for i := 0; i < count; i++ {
		retval = append(retval, server.randomChild(set))
	}
	
	// Also start some new children near the top so that the system doesn't only zoom in
	retval = append(retval, server.randomChild(server.root))

	return retval
}

func multiLevelEvaluation(set *MandelbrotSet) (evaluation float64) {
	set2 := new(MandelbrotSet)

	set2.CenterX = set.CenterX
	set2.CenterY = set.CenterY
	set2.Side = set.Side / 2.0

	evaluation = evaluateBeauty(set, 50)
	evaluation += evaluateBeauty(set2, 60)

	return
}

func evaluateBeauty(set *MandelbrotSet, resolution int) (evaluation float64) {
	bufferSize := resolution
	renderer := new(Renderer)
	maxCount := set.EstimateMaxCount(80)

	bandMap := NewLogarithmicBandMap(maxCount, 35.0)
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
		evaluation *= 1.6

		// But not too much black
		//		r := float64(bufferSize*bufferSize) / float64(pointsInSet) / 1.0
		//		evaluation += r
		//		log.WithFields(log.Fields{
		//			"pointsInSet": pointsInSet,
		//			"blackRatio":  r,
		//			"evaluation":  evaluation,
		//			"maxCount":    maxCount,
		//		}).Info("adjusing beauty")
	}

	// More colors good
	evaluation += float64(h.valueCount()) / 4

	// Less max count good
	evaluation -= float64(maxCount) / 60

	//	log.WithFields(log.Fields{
	//		"colorValues": h.valueCount(),
	//		"evaluation":  evaluation,
	//	}).Info("Returning from evaluateBeauty")
	return evaluation
}

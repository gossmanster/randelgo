package randelbrot

import (
	"math"
)

type bandMap struct {
	values   []int32
	maxCount int
}

// NewLogarithmicBandMap creates a bandMap with a logarithmic distribution of colors to iterations of the Mandelbrot function
func NewLogarithmicBandMap(maxCount int, combinationFactor float64) *bandMap {
	m := bandMap{make([]int32, maxCount), maxCount}

	// combine bands logarithmically
	for i := 0; i < maxCount; i++ {
		temp := math.Log(float64(i)+1) * combinationFactor
		m.values[i] = int32(temp)
	}

	// Now make the bands consecutively numbered so colors are smoother
	var bandNumber int32 = 1
	lastBand := m.values[0]
	for i := 0; i < maxCount; i++ {
		if lastBand != m.values[i] {
			bandNumber++
			lastBand = m.values[i]
		}
		m.values[i] = bandNumber
	}

	return &m
}

func (bandMap *bandMap) Map(count int) int32 {
	if count >= bandMap.maxCount {
		return -1
	}
	if count < 0 {
		return -1
	}
	return bandMap.values[count]
}

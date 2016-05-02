package randelbrot

import (
	"math"
)

// MandelbrotSet represents a location in the imaginary plane that is to be rendered
type MandelbrotSet struct {
	CenterX, CenterY float64
	Side             float64
}

// EstimateMaxCount guesses how many max iterations should be used to render this location
func (m *MandelbrotSet) EstimateMaxCount() int {
	t := math.Log(1.0 / m.Side)
	t = t * 100
	return int(t + 600)
}

// CalculateCount actually calculates the value for a point in the set
func CalculateCount(cx float64, cy float64, maxCount int) int {
	tx := cx
	ty := cy
	x2 := tx * tx
	y2 := ty * ty
	count := 0

	for ((x2 + y2) < 4.0) && (count < maxCount) {
		ty = 2*(tx*ty) + cy
		tx = (x2 - y2) + cx
		x2 = tx * tx
		y2 = ty * ty
		count++
	}
	return count
}

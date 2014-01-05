package randelbrot

import (
	"math"
)

type MandelbrotSet struct {
	CenterX, CenterY float64
	Side             float64
}

func (m *MandelbrotSet) EstimateMaxCount() int {
	t := math.Log(1.0 / m.Side)
	t = t * 80
	return int(t + 600)
}

func CalculateCount(cx float64, cy float64, maxCount int) int {
	tx := cx
	ty := cy
	x2 := tx * tx
	y2 := ty * ty
	count := 0

	for ((x2 + y2) < 4.0) && (count < maxCount) {
		ty = 2*(tx * ty) + cy
		tx = (x2 - y2) + cx
		x2 = tx * tx
		y2 = ty * ty
		count++
	}
	return count
}

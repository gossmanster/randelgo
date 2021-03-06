package randelbrot

import (
	"math"
)

// A Renderer maintains state and provides functions for rendering Mandelbrot Set images
type Renderer struct {
	xCoordinates, yCoordinates []float64
}

func (r *Renderer) initializeCoordinateMap(sizeX int, sizeY int, set *MandelbrotSet, maxCount int) {
	r.xCoordinates = make([]float64, sizeX)
	r.yCoordinates = make([]float64, sizeY)

	size := set.Side
	gapX := size / float64(sizeX)
	gapY := size / float64(sizeY)
	gap := math.Min(gapX, gapY)
	x := set.CenterX - ((gap * float64(sizeX)) / 2.0)
	y := set.CenterY - ((gap * float64(sizeY)) / 2.0)

	for i := 0; i < sizeX; i++ {
		r.xCoordinates[i] = x
		x += gap
	}
	for j := 0; j < sizeY; j++ {
		r.yCoordinates[j] = y
		y += gap
	}
}
func setBand(x int, y int, count int, buffer *PixelBuffer, bandMap *bandMap) {
	band := bandMap.Map(count)
	buffer.SetValue(x, y, band)
}

// Render draws the specified location into the given PixelBuffer
func (r *Renderer) Render(buffer *PixelBuffer, set *MandelbrotSet, bandMap *bandMap, maxCount int) {
	r.initializeCoordinateMap(buffer.SizeX(), buffer.SizeY(), set, maxCount)

	for i := 0; i < buffer.SizeX(); i++ {
		tx := r.xCoordinates[i]
		for j := 0; j < buffer.SizeY(); j++ {
			ty := r.yCoordinates[j]
			count := CalculateCount(tx, ty, maxCount)
			setBand(i, j, count, buffer, bandMap)
		}
	}
}

func (r *Renderer) getOrCalculateBand(buffer *PixelBuffer, bandMap *bandMap, i int, j int, maxCount int) (band int32, calculated bool) {
	calculated = false
	band = 0
	if (i < 0) || (j < 0) || (i >= buffer.SizeX()) || (j >= buffer.SizeY()) {
		return
	}
	band = buffer.GetValue(i, j)
	if band == 0 {
		calculated = true
		count := CalculateCount(r.xCoordinates[i], r.yCoordinates[j], maxCount)
		band = bandMap.Map(count)
		buffer.SetValue(i, j, band)
	}

	return
}

func getBand(buffer *PixelBuffer, i int, j int) (band int32) {
	band = 0
	if (i < 0) || (j < 0) || (i >= buffer.SizeX()) || (j >= buffer.SizeY()) {
		return
	}
	band = buffer.GetValue(i, j)
	return
}

func fillToLeft(buffer *PixelBuffer, i int, j int, band int32) {
	testBand := getBand(buffer, i-1, j)
	if (testBand == 0) || (testBand == band) {
		temp := i - 1
		testBand = getBand(buffer, temp, j)
		for testBand == 0 {
			if temp < 0 {
				return
			}
			if testBand == 0 {
				buffer.SetValue(temp, j, band)
			}
			temp--
			testBand = getBand(buffer, temp, j)
		}
	}
}
func fillCrawl(buffer *PixelBuffer, firstI int, firstJ int, band int32) {
	i := firstI
	j := firstJ
	iinc := 1
	jinc := 1
	done := false
	for !done {
		if getBand(buffer, i+iinc, j) != band {
			jinc = iinc
		} else {
			jinc = -1 * iinc
			i += iinc
			done = ((firstI == i) && (firstJ == j))
			if jinc > 0 {
				fillToLeft(buffer, i, j, band)
			}
		}
		if done {
			break
		}
		if getBand(buffer, i, j+jinc) != band {
			iinc = -1 * jinc
		} else {
			iinc = jinc
			j += jinc
			done = ((firstI == i) && (firstJ == j))
			if jinc > 0 {
				fillToLeft(buffer, i, j, band)
			}
		}
	}
}

func (r *Renderer) crawl(buffer *PixelBuffer, bandMap *bandMap, firstI int, firstJ int, bandInterior int32, maxCount int) (crawled bool) {
	crawled = false
	done := false
	i := firstI
	j := firstJ
	iinc := 1
	jinc := 1
	for !done {
		band, calculated := r.getOrCalculateBand(buffer, bandMap, i+iinc, j, maxCount)
		if band != bandInterior {
			if calculated {
				crawled = true
			}
			jinc = iinc
		} else {
			jinc = -1 * iinc
			i += iinc
			done = ((firstI == i) && (firstJ == j))
		}
		band, calculated = r.getOrCalculateBand(buffer, bandMap, i, j+jinc, maxCount)
		if band != bandInterior {
			if calculated {
				crawled = true
			}
			iinc = -1 * jinc
		} else {
			iinc = jinc
			j += jinc
			done = ((firstI == i) && (firstJ == j))
		}
	}

	return
}

// RenderByCrawling uses the contour crawling algorithm to draw the given location of the set into the PixelBuffer
func (r *Renderer) RenderByCrawling(buffer *PixelBuffer, set *MandelbrotSet, bandMap *bandMap, maxCount int) (numberOfContours int) {
	numberOfContours = 0
	r.initializeCoordinateMap(buffer.SizeX(), buffer.SizeY(), set, maxCount)

	for i := 0; i < buffer.SizeX(); i++ {
		// Keep track of the last band and how many pixels into that band we are
		// Start crawling after we see a few pixels of the same band
		lastBand := int32(0)
		numberOfPointsFoundInBand := 0
		startOfBand := 0
		for j := 0; j < buffer.SizeY(); j++ {
			band, calculated := r.getOrCalculateBand(buffer, bandMap, i, j, maxCount)
			if calculated && (band == lastBand) {
				numberOfPointsFoundInBand++
			} else {
				if band != lastBand {
					startOfBand = j
					lastBand = band
				}
				numberOfPointsFoundInBand = 1
			}
			if numberOfPointsFoundInBand > 5 {
				if r.crawl(buffer, bandMap, i, startOfBand, band, maxCount) {
					numberOfContours++
					fillCrawl(buffer, i, startOfBand, band)
				}
				numberOfPointsFoundInBand = 0
			}
		}
	}
	return
}

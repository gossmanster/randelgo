package utils

import (
	"image/color"
	"math"
)

// Palette maps a set of colors to iteration counts
type Palette struct {
	Colors     []color.RGBA
	ColorCount int
}

var maxCountColor = color.RGBA{0x00, 0x00, 0x00, 0xff}

// NewDefaultPalette creates a Palette and initializes it to an interesting set of values
func NewDefaultPalette() *Palette {

	numberColors := 1024
	retval := Palette{make([]color.RGBA, numberColors), numberColors}

	radianNorm := 9.0 / float64(numberColors)
	for i := 0; i < numberColors; i++ {

		// Completely magic numbers, tune to taste
		// These values are radians which we are going to use to create a Sin curve
		r := float64(i)*(radianNorm*0.05) + 1.0
		b := float64(i) * (radianNorm + 0.11)
		g := float64(i+100)*(radianNorm-0.055) - 2.6

		retval.Colors[i].R = uint8(math.Abs(math.Sin(r))*239.0) + 16
		retval.Colors[i].G = uint8(math.Abs(math.Sin(g) * 255.0))
		retval.Colors[i].B = uint8(math.Abs(math.Sin(b)) * 255.0)

		retval.Colors[i].A = uint8(255)
	}

	return &retval
}

// Map takes an interation count and returns the color it should use
func (p *Palette) Map(value int) color.RGBA {
	if value == -1 {
		return maxCountColor
	}
	return p.Colors[value%p.ColorCount]
}

package utils

import (
    "image/color"
    "math"
)

type Palette struct {
    Colors     []color.RGBA
    ColorCount int
}

var MaxCountColor = color.RGBA{ 0x00, 0x00, 0x00, 0xff }

func NewDefaultPalette() *Palette {

    numberColors := 1024
    retval := Palette { make([]color.RGBA, numberColors), numberColors}

    radianNorm := 9.0 / float64(numberColors)
    for i := 0; i < numberColors; i++ {

        // Completely magic numbers, tune to taste
        // These values are radians which we are going to use to create a Sin curve
        r := float64(i)*(radianNorm*0.05) + 1.0
        g := float64(i) * (radianNorm + 0.11)
        b := float64(i+100)*(radianNorm-0.055) - 2.6

        retval.Colors[i].B = uint8(math.Abs(math.Sin(b)) * 256.0)
        retval.Colors[i].R = uint8(math.Abs(math.Sin(r))*256.0) + 16
        retval.Colors[i].G = uint8(math.Abs(math.Sin(g) * 256.0))
        retval.Colors[i].A = uint8(255)
    }

    return &retval
}

func (self *Palette) Map(value int) color.RGBA {
    if (value == -1) {
        return MaxCountColor
    }
    return self.Colors[value % self.ColorCount]
}

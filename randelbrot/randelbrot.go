package randelbrot

import (
	"image"
	"image/color"
)

func DrawImage() image.Image {
	m := image.NewRGBA(image.Rect(0, 0, 400, 400))
	blue := color.RGBA{0, 255, 255, 255}

	for i := 0; i < 400; i++ {
		m.SetRGBA(i, i, blue)
	}

	return m
}

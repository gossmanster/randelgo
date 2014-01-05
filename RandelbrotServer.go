package main

import (
    "github.com/user/randelbrot/randelbrot"
    "github.com/user/randelbrot/utils"
    "image"
    "image/jpeg"
    "net/http"
)

func main() {
    http.HandleFunc("/", foo)
    http.ListenAndServe(":1966", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "image/jpeg")

    m := render()

    jpeg.Encode(w, m, nil)
}

func render() *image.RGBA {
    set := randelbrot.MandelbrotSet{-0.75, 0.0, 2.5}
    maxCount := set.EstimateMaxCount()
    buffer := randelbrot.NewPixelBuffer(1000, 1000)
    bandMap := randelbrot.NewLogarithmicBandMap(maxCount, 32.0)

    renderer := new(randelbrot.Renderer)
    renderer.RenderByCrawling(buffer, &set, bandMap, maxCount)

    return convertToImage(buffer)
}

func convertToImage(buffer *randelbrot.PixelBuffer) *image.RGBA {

    palette := utils.NewDefaultPalette()
    r := image.Rect(0, 0, buffer.SizeX(), buffer.SizeY())
    img := image.NewRGBA(r)

    for x := 0; x < buffer.SizeX(); x++ {
        for y := 0; y < buffer.SizeY(); y++ {
            count := buffer.GetValue(x, y)
            img.SetRGBA(x, y, palette.Map(int(count)))
        }
    }

    return img
}

package main

import (
    "randelgo/randelbrot"
    "randelgo/utils"
    "image"
    "image/jpeg"
    "net/http"
    "runtime"
)

func main() {

    runtime.GOMAXPROCS(2)
    renderChannel := make(chan *image.RGBA, 10)
    go render(renderChannel)

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "image/jpeg")

            m := <-renderChannel

            jpeg.Encode(w, m, nil)
        })
    http.ListenAndServe(":1966", nil)
}



func render(renderChannel chan *image.RGBA) {
    renderer := new(randelbrot.Renderer)

    x := -1.6

    for { 
        set := randelbrot.MandelbrotSet{x, 0.0, 0.5}
        maxCount := set.EstimateMaxCount()
        buffer := randelbrot.NewPixelBuffer(1000, 1000)
        bandMap := randelbrot.NewLogarithmicBandMap(maxCount, 32.0)


        renderer.RenderByCrawling(buffer, &set, bandMap, maxCount)

        renderChannel <- convertToImage(buffer) 
        x += 0.1
    }    
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

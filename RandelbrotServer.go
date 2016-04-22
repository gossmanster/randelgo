package main

import (
	log "github.com/Sirupsen/logrus"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"
	"net/http"
	"randelgo/randelbrot"
	"randelgo/utils"
	"time"
)

func main() {
	log.Info("Starting Randelgo Server")
	renderChannel := make(chan *image.RGBA, 30)
	go render(renderChannel)
	
	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Automatic Mandelbrot Explorer\n")
		fmt.Fprintln(w, "/newImage to get a JPG format image")
	})

	http.HandleFunc("/newImage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")

		m := <-renderChannel

		jpeg.Encode(w, m, nil)
	})
	http.ListenAndServe(":80", nil)
}

func render(renderChannel chan *image.RGBA) {
	n := int64(time.Now().Nanosecond())

	r := rand.New(rand.NewSource(n))

	server := randelbrot.NewRandelbrotServer(r)

	for {

		buffer := randelbrot.NewPixelBuffer(1000, 1000)
		server.RenderNext(buffer)

		renderChannel <- convertToImage(buffer)
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

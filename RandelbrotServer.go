package main

import (
	log "github.com/Sirupsen/logrus"
	"image"
	"image/jpeg"
	"math/rand"
	"net/http"
	"randelgo/randelbrot"
	"randelgo/utils"
	"runtime"
	"time"
)

func main() {
	log.Info("Starting Randelgo Server")
	runtime.GOMAXPROCS(2)
	renderChannel := make(chan *image.RGBA, 30)
	go render(renderChannel)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")

		m := <-renderChannel

		jpeg.Encode(w, m, nil)
	})
	http.ListenAndServe(":1966", nil)
}

func render(renderChannel chan *image.RGBA) {
	n := int64(time.Now().Nanosecond())

	r := rand.New(rand.NewSource(89))
	rand.Seed(n)
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

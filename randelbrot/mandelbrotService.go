package randelbrot

func RenderToBuffer(buffer *PixelBuffer, set *MandelbrotSet) {
    renderer := new(Renderer)
    maxCount := set.EstimateMaxCount()

	bandMap := NewLogarithmicBandMap(maxCount, 32.0)  

	renderer.RenderByCrawling(buffer, set, bandMap, maxCount) 
}
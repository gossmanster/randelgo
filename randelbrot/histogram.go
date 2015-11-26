package randelbrot

type histogram struct {
	values map[int32]int
}

func evaluateBuffer(buffer *PixelBuffer) *histogram {
	h := new(histogram)
	h.values = make(map[int32]int)

	for x := 0; x < buffer.SizeX(); x++ {
		for y := 0; y < buffer.SizeY(); y++ {
			count := buffer.GetValue(x, y)
			h.values[count] = h.values[count] + 1
		}
	}

	return h
}

func (h *histogram) valueCount() int {
	return len(h.values)
}

func (h *histogram) getValue(count int32) int {
	return h.values[count]
}

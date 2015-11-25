package randelbrot

// PixelBuffer is the backing plane for almost any kind of rendering
type PixelBuffer struct {
	bits                                   []int32
	offsetX, offsetY, sizeX, sizeY, stride int
}

// NewPixelBuffer initializes a PixelBuffer of the specified size
func NewPixelBuffer(sizeX int, sizeY int) *PixelBuffer {
	bits := make([]int32, sizeX*sizeY)
	retval := PixelBuffer{bits, 0, 0, sizeX, sizeY, sizeX}
	return &retval
}

// SizeX is a getter for the X dimension of the buffer
func (buffer *PixelBuffer) SizeX() int {
	return buffer.sizeX
}

// SizeY is a getter for the Y dimension of the buffer
func (buffer *PixelBuffer) SizeY() int {
	return buffer.sizeY
}

// Clear sets all the pixels to 0 (black depending on how the final renderer works)
func (buffer *PixelBuffer) Clear() {
	for i := 0; i < len(buffer.bits); i++ {
		buffer.bits[i] = 0
	}

}

func (buffer *PixelBuffer) index(x int, y int) int {
	offset := (x + buffer.offsetX) + ((y + buffer.offsetY) * buffer.stride)
	return offset
}

// GetValue does just that for a specific pixel
func (buffer *PixelBuffer) GetValue(x int, y int) int32 {
	return buffer.bits[buffer.index(x, y)]
}

// SetValue does just that for a specific pixel
func (buffer *PixelBuffer) SetValue(x int, y int, value int32) {
	buffer.bits[buffer.index(x, y)] = value
}

// GetPixels returns the underlying array. Used for example to convert to a bitmap
func (buffer *PixelBuffer) GetPixels() []int32 {
	return buffer.bits
}

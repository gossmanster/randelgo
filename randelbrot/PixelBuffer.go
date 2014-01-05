package randelbrot

type PixelBuffer struct {
	bits []int32
	offsetX, offsetY, sizeX, sizeY, stride int
}

func NewPixelBuffer(sizeX int, sizeY int) *PixelBuffer {
	bits := make([]int32, sizeX*sizeY)
	retval := PixelBuffer{bits, 0, 0, sizeX, sizeY, sizeX}
	return &retval
}

func (self *PixelBuffer) SizeX() int {
	return self.sizeX
}

func (self *PixelBuffer) SizeY() int {
	return self.sizeY
}

func (self *PixelBuffer) Clear() {
	for i := 0; i < len(self.bits); i++ {
		self.bits[i] = 0
	}

}

func (self *PixelBuffer) index(x int, y int) int {
	offset := (x + self.offsetX) + ((y + self.offsetY) * self.stride)
	return offset
}

func (self *PixelBuffer) GetValue(x int, y int) int32 {
	return self.bits[self.index(x, y)]
}

func (self *PixelBuffer) SetValue(x int, y int, value int32) {
	self.bits[self.index(x, y)] = value
}

func (self *PixelBuffer) GetPixels() []int32 {
	return self.bits
}

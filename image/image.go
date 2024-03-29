package main

import (
	"golang.org/x/tour/pic"
	"image"
	"image/color"
)

type Image struct{
	w, h int
}

func (img Image) ColorModel() color.Model {
	return color.RGBAModel
}
func (img Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.w, img.h)
}
func (img Image) At(x, y int) color.Color {
	return color.RGBA{uint8(x+y), uint8(x+y), 255, 255}
}

func main() {
	m := Image{200, 200}
	pic.ShowImage(m)
}
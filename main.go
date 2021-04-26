package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"

	"syscall/js"
)

func register() {
	js.Global().Set("convertToGray", js.FuncOf(convertToGray))
}
func main() {
	fin := make(chan struct{}, 0)
	fmt.Println("Module WASM chargé!")
	register()

	<-fin
}
func convertToGray(this js.Value, args []js.Value) interface{} {

	file := bytes.NewReader(typedArrayToByteSlice(args[0]))
	oldimage, err := jpeg.Decode(file)
	if err != nil {
		log.Print(err.Error())

	}
	b := oldimage.Bounds()

	l, h := b.Max.X, b.Max.Y

	gsimage := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{l, h}})
	for i := 0; i < l; i++ {
		for j := 0; j < h; j++ {
			pixelcolor := oldimage.At(i, j)
			red, green, blue, _ := pixelcolor.RGBA()
			r := math.Pow(float64(red), 2.6)
			g := math.Pow(float64(green), 2.6)
			b := math.Pow(float64(blue), 2.6)
			m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.6)
			Y := uint16(m + 0.5)
			gc := color.Gray{uint8(Y >> 8)}
			gsimage.Set(i, j, gc)

		}
	}
	newfile := new(bytes.Buffer)
	jpeg.Encode(newfile, gsimage, nil)
	dst := js.Global().Get("Uint8Array").New(len(newfile.Bytes()))

	js.CopyBytesToJS(dst, newfile.Bytes())

	return dst
}
func typedArrayToByteSlice(arg js.Value) []byte {
	length := arg.Length()
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(arg.Index(i).Int())
	}
	return bytes
}

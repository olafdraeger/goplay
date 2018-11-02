package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
)

const colors = 64

const (
	whiteIndex = 0
	blackIndex = 1
)

func lissajous(out io.Writer, cyclesRequested float64) {

	var palette = []color.Color{color.Black}
	for j := 0; j < colors; j++ {
		r := uint8(rand.Intn(256))
		g := uint8(rand.Intn(256))
		b := uint8(rand.Intn(256))

		palette = append(palette, color.RGBA{r, g, b, 1})
	}
	const (
		res     = 0.001
		size    = 100
		nframes = 64
		delay   = 8
	)
	var cycles float64
	if cyclesRequested != 0 {
		cycles = cyclesRequested
	} else {
		cycles = 5
	}
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			//paletteIndex := uint8(rand.Intn(64))
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(i))
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}

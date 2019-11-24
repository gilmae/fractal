package fractal

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"strconv"

	"github.com/gilmae/interpolation"
)

var xSequence []float64
var redpoints []float64
var greenpoints []float64
var bluepoints []float64

var redInterpolant interpolation.MonotonicCubic
var greenInterpolant interpolation.MonotonicCubic
var blueInterpolant interpolation.MonotonicCubic

var palette = make([]color.NRGBA, paletteLength)
var paletteLength int = 16

func build_gradient(gradient_str string) {
	var g [][]string

	byt := []byte(gradient_str)
	_ = json.Unmarshal(byt, &g)
	var size = len(g)
	xSequence = make([]float64, size)
	redpoints = make([]float64, size)
	greenpoints = make([]float64, size)
	bluepoints = make([]float64, size)

	for i, v := range g {
		xSequence[i], _ = strconv.ParseFloat(v[0], 64)
		b, _ := hex.DecodeString(v[1])
		redpoints[i] = float64(b[0])
		greenpoints[i] = float64(b[1])
		bluepoints[i] = float64(b[2])
	}

	redInterpolant = interpolation.CreateMonotonicCubic(xSequence, redpoints)
	greenInterpolant = interpolation.CreateMonotonicCubic(xSequence, greenpoints)
	blueInterpolant = interpolation.CreateMonotonicCubic(xSequence, bluepoints)
}

func fill_palette() {

	for i := 0; i < paletteLength; i++ {
		var point = 1.0 * float64(i) / float64(paletteLength)
		var redpoint = redInterpolant(point)
		var greenpoint = greenInterpolant(point)
		var bluepoint = blueInterpolant(point)

		palette[i] = color.NRGBA{uint8(redpoint), uint8(greenpoint), uint8(bluepoint), 255}

	}
}

func get_colour(esc int, maxIterations int, colour_mode string) color.NRGBA {
	var escapeAsFloat = float64(esc)

	if colour_mode == "true" {

		var point = escapeAsFloat / float64(maxIterations)
		var redpoint = redInterpolant(point)
		var greenpoint = greenInterpolant(point)
		var bluepoint = blueInterpolant(point)

		return color.NRGBA{uint8(redpoint), uint8(greenpoint), uint8(bluepoint), 255}
	} else if colour_mode == "smooth" {
		index1 := int(math.Abs(escapeAsFloat))
		t2 := escapeAsFloat - float64(index1)
		t1 := 1 - t2

		index1 = index1 % len(palette)
		index2 := (index1 + 1) % len(palette)

		clr1 := palette[index1]
		clr2 := palette[index2]

		r := float64(clr1.R)*t1 + float64(clr2.R)*t2
		g := float64(clr1.G)*t1 + float64(clr2.G)*t2
		b := float64(clr1.B)*t1 + float64(clr2.B)*t2

		return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
	} else if colour_mode == "banded" {
		return palette[int(esc)%len(palette)]
	} else {
		return color.NRGBA{255, 255, 255, 255}
	}
}

func Draw_Image(filename string, plot_map map[Key]Point, width int, height int, gradient string, maxIterations int, colour_mode string) {
	build_gradient(gradient)
	fill_palette()

	bounds := image.Rect(0, 0, width, height)

	b := image.NewNRGBA(bounds)
	draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var p = plot_map[Key{x, y}]
			if !p.Escaped {
				b.Set(p.X, p.Y, color.NRGBA{0, 0, 0, 255})

			} else {
				b.Set(p.X, p.Y, get_colour(p.Escape, maxIterations, colour_mode))
			}
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}

	if err = jpeg.Encode(file, b, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
		fmt.Println(err)
	}

	if err = file.Close(); err != nil {
		fmt.Println(err)
	}
}

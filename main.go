package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func toDigits(v int) []int {
	digits := make([]int, 0, 10)
	for {
		d := v % 10
		digits = append(digits, d)
		v /= 10
		if v == 0 {
			break
		}
	}
	return digits
}

func loadPng(filename string) image.Image {
	f, err := os.Open(filename)
	check(err)
	defer f.Close()

	img, err := png.Decode(f)
	check(err)

	return img
}

func savePng(filename string, img image.Image) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()

	check(png.Encode(f, img))
}

type template struct {
	digits     []image.Image
	background image.Image
}

func loadImages() *template {
	r := template{
		digits:     make([]image.Image, 10),
		background: loadPng("template/background.png"),
	}
	for i := 0; i < 10; i++ {
		r.digits[i] = loadPng(fmt.Sprintf("template/%d.png", i))
	}
	return &r
}

func generateNumberImage(v int, t *template) image.Image {
	digits := toDigits(v)

	n := len(digits)
	dr := t.digits[8].Bounds()
	ir := t.background.Bounds()

	p := image.Point{
		X: (ir.Dx() - dr.Dx()*n) / 2,
		Y: (ir.Dy() - dr.Dy()) / 2,
	}

	m := image.NewRGBA(ir.Bounds())
	draw.Draw(m, m.Bounds(), t.background, image.ZP, draw.Src)

	for i := len(digits) - 1; i >= 0; i-- {
		dst := image.Rectangle{Min: p, Max: p.Add(dr.Size())}
		draw.Draw(m, dst, t.digits[digits[i]], image.ZP, draw.Over)
		p.X += dr.Dx()
	}

	return m
}

func ensureDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func main() {
	t := loadImages()
	const outputDirectory = "numbers"
	ensureDirectory(outputDirectory)

	for i := 1; i < 1000; i++ {
		image := generateNumberImage(i, t)
		savePng(fmt.Sprintf("%s/%03d.png", outputDirectory, i), image)
	}
}

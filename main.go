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

func toDigits(v int) ([]int) {
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

func loadPng(filename string) (image.Image) {
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

func loadImages() ([]image.Image, image.Image) {
	images := make([]image.Image, 10)
	for i := 0; i < 10; i++ {
		images[i] = loadPng(fmt.Sprintf("template/%d.png", i))
	}
	background := loadPng("template/background.png")
	return images, background
}


func generateNumberImage(v int, images []image.Image, background image.Image) (image.Image) {
	digits := toDigits(v)

	n := len(digits)
	dr := images[8].Bounds()
	ir := background.Bounds()

	p := image.Point{
		X: (ir.Dx() - dr.Dx() * n) / 2,
		Y: (ir.Dy() - dr.Dy()) / 2,
	}

	m := image.NewRGBA(ir.Bounds())
	draw.Draw(m, m.Bounds(), background, image.ZP, draw.Src)

	for i := len(digits) - 1; i >= 0; i-- {
		dst := image.Rectangle{Min: p, Max: p.Add(dr.Size())}
		draw.Draw(m, dst, images[digits[i]], image.ZP, draw.Over)
		p.X += dr.Dx()
	}

	return m
}

func main() {
	images, background := loadImages()
	_ = os.Mkdir("numbers", os.ModePerm)

	for i := 1; i < 1000; i++ {
		image := generateNumberImage(i, images, background)
		savePng(fmt.Sprintf("numbers/%03d.png", i), image)
	}
}

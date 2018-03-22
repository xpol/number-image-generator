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

type theme struct {
	digits      []image.Image
	backgrounds []image.Image
}

func loadTheme(themeName string) *theme {
	r := theme{
		digits:      make([]image.Image, 10),
		backgrounds: make([]image.Image, 10),
	}

	for i := 0; i < 10; i++ {
		r.digits[i] = loadPng(fmt.Sprintf("themes/%s/digits/%d.png", themeName, i))
	}

	for i := 0; i < 10; i++ {
		r.backgrounds[i] = loadPng(fmt.Sprintf("themes/%s/backgrounds/%d.png", themeName, i))
	}

	return &r
}

func generateNumberImage(v int, t *theme) image.Image {
	digits := toDigits(v)

	n := len(digits)
	dr := t.digits[8].Bounds()
	ir := t.backgrounds[0].Bounds()

	p := image.Point{
		X: (ir.Dx() - dr.Dx()*n) / 2,
		Y: (ir.Dy() - dr.Dy()) / 2,
	}

	m := image.NewRGBA(ir.Bounds())
	draw.Draw(m, m.Bounds(), t.backgrounds[(v-1)/10%10], image.ZP, draw.Src)

	for i := len(digits) - 1; i >= 0; i-- {
		dst := image.Rectangle{Min: p, Max: p.Add(dr.Size())}
		draw.Draw(m, dst, t.digits[digits[i]], image.ZP, draw.Over)
		p.X += dr.Dx()
	}

	return m
}

func ensureDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

func main() {
	themeName := "default"
	if len(os.Args) >= 2 {
		themeName = os.Args[1]
	}
	t := loadTheme(themeName)
	outputDirectory := fmt.Sprintf("numbers/%s", themeName)
	ensureDirectory(outputDirectory)

	for i := 1; i < 1000; i++ {
		image := generateNumberImage(i, t)
		savePng(fmt.Sprintf("%s/%03d.png", outputDirectory, i), image)
	}
}

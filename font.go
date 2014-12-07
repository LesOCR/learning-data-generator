package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"regexp"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"code.google.com/p/go.image/bmp"
)

func main() {
	// Font loading
	fonts := []*truetype.Font{
		// Sans-serif
		readFont("arial.ttf"),
		readFont("verdana.ttf"),
		readFont("trebuchet.ttf"),
		readFont("microsoft-sans.ttf"),
		readFont("merriweather-sans.ttf"),

		// Serif
		readFont("times-new-roman.ttf"),
		readFont("georgia.ttf"),
	}

	reg := regexp.MustCompile("[A-Za-z0-9,.?!\\-_]")

	for i := byte(0); i < 255; i++ {
		if !reg.Match([]byte{i}) {
			continue
		}
		for j, font := range fonts {
			img, err := MakeImage(string(i), font)
			if err != nil {
				panic(err)
			}
			os.Mkdir(fmt.Sprintf("./out/%d", i), 0777)
			saveToBmpFile(fmt.Sprintf("out/%d/%d.bmp", i, j), img)
		}
	}
}

// MakeImage generates an image with the specified text at the specified size
// (in bold type or not)
func MakeImage(text string, font *truetype.Font) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, 35, 35))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFontSize(float64(30))
	c.SetDst(img)
	c.SetClip(img.Bounds())
	c.SetSrc(image.NewUniform(color.Black))
	c.SetFont(font)

	tw := textWidth(font, 30, text)
	pt := freetype.Pt(15-int(tw/2), 26)
	_, err := c.DrawString(text, pt)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// saveToBmpFile saves a BMP to a file
func saveToBmpFile(filePath string, i *image.RGBA) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = bmp.Encode(f, i)
	if err != nil {
		panic(err)
	}
}

// readFont opens a font file and makes freetype decode it
func readFont(file string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	return font
}

// textWidth measures very accurately the width a text uses on an image
func textWidth(f *truetype.Font, size int, s string) float64 {
	scale, width := float64(size)/float64(f.FUnitsPerEm()), 0
	prev, hasPrev := truetype.Index(0), false
	for _, rune := range s {
		index := f.Index(rune)
		if hasPrev {
			width += int(f.Kerning(f.FUnitsPerEm(), prev, index))
		}
		width += int(f.HMetric(f.FUnitsPerEm(), index).AdvanceWidth)
		prev, hasPrev = index, true
	}
	return float64(width) * scale
}

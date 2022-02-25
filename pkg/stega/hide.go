package stega

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
)

type PixelColors struct {
	Colors []color.RGBA64
}

func (pc *PixelColors) NullifyLSB() {
	for i := range pc.Colors {
		pc.Colors[i].B &= ^uint16(1)
	}
}

func HideInfo(originalFileName, secretText, result string) {
	image.RegisterFormat("bmp", "bmp", bmp.Decode, bmp.DecodeConfig)

	f, err := os.Open(originalFileName)
	if err != nil {
		log.Fatalln("error:", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalln(err)
	}

	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	prepImg := NewPixelColorsFromImage(img, width, height)
	prepImg.NullifyLSB()

	sb := strings.Builder{}
	for _, b := range []byte(secretText) {
		sb.WriteString(fmt.Sprintf("%.08b", b))
	}

	for i, sym := range sb.String() {
		if sym == '1' {
			prepImg.Colors[i].B |= 1
		}
	}

	out := image.NewNRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})
	var k int
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			out.Set(x, y, prepImg.Colors[k])
			k++
		}
	}

	outf, _ := os.Create(result)
	defer outf.Close()
	bmp.Encode(outf, out)
}

func NewPixelColorsFromImage(img image.Image, width, height int) (pc PixelColors) {
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			pc.Colors = append(pc.Colors, color.RGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(a),
			})
		}
	}
	return
}

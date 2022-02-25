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

func init() {
	image.RegisterFormat("bmp", "bmp", bmp.Decode, bmp.DecodeConfig)
}

type PixelColors struct {
	Colors []color.NRGBA
}

func (pc *PixelColors) NullifyLSB(bits int) {
	for i := 0; i < bits; i++ {
		pc.Colors[i].B &= ^uint8(1)
	}
}

func HideInfo(originalFileName, secretText, result string) {
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
	prepImg.NullifyLSB(len(secretText) * 16)

	sb := strings.Builder{}
	for _, b := range []rune(secretText) {
		sb.WriteString(fmt.Sprintf("%.016b", b))
	}

	var n int
	for _, sym := range sb.String() {
		if n == len(secretText)*16 {
			break
		}
		if sym == '1' {
			prepImg.Colors[n].B += 1
		}
		n++
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
			pc.Colors = append(pc.Colors, color.NRGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
		}
	}
	return
}

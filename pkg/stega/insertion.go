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
	Colors []color.RGBA
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

	if len(secretText)*16 > width*height*3 {
		log.Fatalln("Встраивание ЦВЗ в стеганоконтейнер невозможно")
	}

	pc := NewPixelColorsFromImage(img, width, height)

	fmt.Println("Значения синей цветовой компоненты в диапазоне длины сообщения до встраивания")
	for i := 0; i <= len(secretText); i++ {
		fmt.Printf("%v ", pc.Colors[i].B)
	}
	fmt.Println()
	pc.NullifyLSB(len(secretText) * 16)

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
			pc.Colors[n].B++
		}
		n++
	}
	fmt.Println("Значения синей цветовой компоненты в диапазоне длины сообщения после встраивания")
	for i := 0; i <= len(secretText); i++ {
		fmt.Printf("%v ", pc.Colors[i].B)
	}
	fmt.Println()
	out := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})

	fo, _ := os.Create("after.txt")
	defer fo.Close()
	var k int
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			out.Set(x, y, pc.Colors[k])
			fmt.Fprintf(fo, "Pixel(%d,%d)={%d,%d,%d}", x, y, pc.Colors[k].R, pc.Colors[k].G, pc.Colors[k].B)
			k++
		}
		fmt.Fprintln(fo)
	}

	outf, _ := os.Create(result)
	defer outf.Close()
	bmp.Encode(outf, out)
}

func NewPixelColorsFromImage(img image.Image, width, height int) (pc PixelColors) {
	f, _ := os.Create("before.txt")
	defer f.Close()
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			pc.Colors = append(pc.Colors, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
			fmt.Fprintf(f, "Pixel(%d,%d)={%d,%d,%d}", i, j, uint8(r), uint8(g), uint8(b))
		}
		fmt.Fprintln(f)
	}
	return
}

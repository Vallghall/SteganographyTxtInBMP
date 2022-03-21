package stega

import (
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
	"unicode/utf8"
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

	if utf8.RuneCountInString(secretText)*16 > width*height {
		log.Fatalln("Встраивание ЦВЗ в стеганоконтейнер невозможно")
	}

	pc := NewPixelColorsFromImage(img, width, height)

	fmt.Println("Значения компонент синего цвета фрагмента изображения до встраивания:")
	for i := 0; i <= 8; i++ {
		for j := 0; j < 16; j++ {
			fmt.Printf("%.08b ", pc.Colors[i*width+j].B)
		}
		fmt.Println("")
	}

	pc.NullifyLSB(utf8.RuneCountInString(secretText) * 16)

	sb := strings.Builder{}
	for _, b := range []rune(secretText) {
		sb.WriteString(fmt.Sprintf("%.016b", b))
	}

	var n int
	for _, sym := range sb.String() {

		if n == utf8.RuneCountInString(secretText)*16 {
			break
		}
		if sym == '1' {
			pc.Colors[n].B++
		}
		n++
	}
	fmt.Println("Значения компоненты синего цвета фрагмента изображения после встраивания:")
	for i := 0; i <= 8; i++ {
		for j := 0; j < 16; j++ {
			fmt.Printf("%.08b ", pc.Colors[i*width+j].B)
		}
		fmt.Println("")
	}

	out := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})

	var k int
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			out.Set(x, y, pc.Colors[k])
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
			pc.Colors = append(pc.Colors, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			})
		}
	}
	return
}

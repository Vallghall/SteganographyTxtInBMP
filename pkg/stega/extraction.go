package stega

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func ExtractLSBInfo(infoLength int, original, result string) string {
	originalF, _ := os.Open(original)
	defer originalF.Close()

	resultF, _ := os.Open(result)
	defer resultF.Close()

	originalImg, _, err := image.Decode(originalF)
	if err != nil {
		log.Fatalln(err)
	}

	resultImg, _, err := image.Decode(resultF)
	if err != nil {
		log.Fatalln(err)
	}

	oWidth, oHeight := originalImg.Bounds().Dx(), originalImg.Bounds().Dy()
	rWidth, rHeight := resultImg.Bounds().Dx(), resultImg.Bounds().Dy()

	if oWidth != rWidth || oHeight != rHeight {
		log.Fatalln("Images are not equal")
	}

	pc1 := NewPixelColorsFromImage(originalImg, oWidth, oHeight)
	pc2 := NewPixelColorsFromImage(resultImg, rWidth, rHeight)

	fmt.Printf("СКО: %.04f\n", MeanSquareError(pc1, pc2))
	fmt.Printf("НСКО: %f\n", NormalizedMeanSquareError(pc1, pc2))

	pc1.NullifyLSB(infoLength * 16)
	bs := extractSecretInfoBitString(infoLength*16, pc1, pc2)

	return extractInfo(bs)
}

func extractSecretInfoBitString(limit int, pc1, pc2 PixelColors) string {
	var sb strings.Builder
	for i := 0; i < limit; i++ {
		if pc1.Colors[i] == pc2.Colors[i] {
			sb.WriteRune('0')
		} else {
			sb.WriteRune('1')
		}
	}
	return sb.String()
}

func extractInfo(bitString string) string {
	var sb strings.Builder

	for i := 0; i < len(bitString)/16; i++ {
		sub := bitString[i*16 : i*16+16]
		intSub, _ := strconv.ParseUint(sub, 2, 32)
		r := rune(intSub)
		sb.WriteRune(r)
	}

	return sb.String()
}

func MeanSquareError(pc1, pc2 PixelColors) float64 {
	var up, length int

	for i, origin := range pc1.Colors {
		length++
		a := int(origin.B)
		b := int(pc2.Colors[i].B)
		up += (a - b) * (a - b)
	}

	return math.Sqrt(float64(up) / float64(length*3))
}

func NormalizedMeanSquareError(pc1, pc2 PixelColors) float64 {
	var up float64
	var length int

	for _, origin := range pc2.Colors {
		length++
		up += math.Pow(float64(origin.B), 2)
	}

	return MeanSquareError(pc1, pc2) / math.Sqrt(up/float64(length*3))
}

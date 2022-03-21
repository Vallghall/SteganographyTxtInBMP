package stega

import (
	"image"
	"log"
	"math"
	"os"
)

func EvalQuality(original, result string) (float64, float64) {
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

	width, height := originalImg.Bounds().Dx(), originalImg.Bounds().Dy()

	pc1 := NewPixelColorsFromImage(originalImg, width, height)
	pc2 := NewPixelColorsFromImage(resultImg, width, height)

	return MeanSquareError(pc1, pc2, width, height), NormalizedMeanSquareError(pc1, pc2, width, height)
}

func MeanSquareError(pc1, pc2 PixelColors, x, y int) float64 {
	var length int
	var up float64

	for i, origin := range pc1.Colors {
		length++
		c := int(origin.B)
		s := int(pc2.Colors[i].B)
		up += math.Pow(float64(c-s), 2)
	}

	return math.Sqrt(up) / float64(x*y)
}

func NormalizedMeanSquareError(pc1, pc2 PixelColors, x, y int) float64 {
	var up, div float64
	var length int

	for i, origin := range pc1.Colors {
		length++
		c := int(origin.B)
		s := int(pc2.Colors[i].B)
		up += math.Pow(float64(c-s), 2)
		div += math.Pow(float64(c), 2)
	}

	return up / div
}

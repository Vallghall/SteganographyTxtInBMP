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

	return MeanSquareError(pc1, pc2), NormalizedMeanSquareError(pc1, pc2)
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

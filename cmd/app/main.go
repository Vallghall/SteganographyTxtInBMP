package main

import (
	"fmt"
	"stega/pkg/stega"
)

const (
	originalFileName = "oberon.bmp"
	secretText       = "Гусев Роман Михайлович"
	result           = "encoded.bmp"
)

func main() {
	stega.HideInfo(originalFileName, secretText, result)

	secretGot := stega.ExtractLSBInfo(len(secretText), originalFileName, result)
	fmt.Println(secretGot)

	mse, nmse := stega.EvalQuality(originalFileName, result)
	fmt.Printf("СКО: %.04f\nНСКО: %f\n", mse, nmse)
}

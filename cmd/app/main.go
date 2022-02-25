package main

import (
	"fmt"
	"stega/pkg/stega"
)

const (
	originalFileName = "oberon.bmp"
	secretText       = "Vali"
	result           = "encoded.bmp"
)

func main() {
	stega.HideInfo(originalFileName, secretText, result)

	fmt.Println(stega.ExtractLSBInfo(len(secretText), originalFileName, result))
}

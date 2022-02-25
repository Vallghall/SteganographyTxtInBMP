package main

import (
	"fmt"
	"stega/pkg/stega"
)

const (
	originalFileName = "oberon.bmp"
	secretText       = "Русский"
	result           = "encoded.bmp"
)

func main() {
	stega.HideInfo(originalFileName, secretText, result)

	fmt.Println(stega.ExtractLSBInfo(len(secretText)*16, originalFileName, result))
}

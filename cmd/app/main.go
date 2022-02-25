package main

import "stega/pkg/stega"

const (
	originalFileName = "oberon.bmp"
	secretText       = "ГусевРоманМихайлович"
	result           = "encoded.bmp"
)

func main() {
	stega.HideInfo(originalFileName, secretText, result)

}

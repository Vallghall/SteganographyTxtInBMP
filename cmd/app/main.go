package main

import (
	"fmt"
	"os"
	"stega/pkg/stega"
)

const (
	originalFileName = "oberon.bmp"
	secretFileName   = "secret.txt"
	result           = "encoded.bmp"
)

func main() {
	secretText, _ := os.ReadFile(secretFileName) // Чтение текстового файла

	stega.HideInfo(originalFileName, string(secretText), result)

	secretGot := stega.ExtractLSBInfo(string(secretText), result)
	fmt.Printf("Извлеченное сообщение: %s\n", secretGot)

	mse, nmse := stega.EvalQuality(originalFileName, result)
	fmt.Printf("СКО: %.08f\nНСКО: %.020f\n", mse, nmse)
}

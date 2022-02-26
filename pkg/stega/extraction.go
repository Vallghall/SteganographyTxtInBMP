package stega

import (
	"image"
	"log"
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

	width, height := originalImg.Bounds().Dx(), originalImg.Bounds().Dy()

	pc := NewPixelColorsFromImage(resultImg, width, height)
	bs := extractSecretInfoBitString(infoLength*16, pc)

	return extractInfo(bs)
}

func extractSecretInfoBitString(limit int, pc PixelColors) string {
	var sb strings.Builder
	for i := 0; i < limit; i++ {
		if c := pc.Colors[i].B; c == c & ^uint8(1) {
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

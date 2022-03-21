package stega

import (
	"image"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func ExtractLSBInfo(info string, result string) string {
	resultF, _ := os.Open(result)
	defer resultF.Close()

	resultImg, _, err := image.Decode(resultF)
	if err != nil {
		log.Fatalln(err)
	}

	width, height := resultImg.Bounds().Dx(), resultImg.Bounds().Dy()

	pc := NewPixelColorsFromImage(resultImg, width, height)
	bs := extractSecretInfoBitString(utf8.RuneCountInString(info)*16, pc)

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

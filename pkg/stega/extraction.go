package stega

import (
	"bytes"
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

	oWidth, oHeight := originalImg.Bounds().Dx(), originalImg.Bounds().Dy()
	rWidth, rHeight := resultImg.Bounds().Dx(), resultImg.Bounds().Dy()

	if oWidth != rWidth || oHeight != rHeight {
		log.Fatalln("Images are not equal")
	}

	pc1 := NewPixelColorsFromImage(originalImg, oWidth, oHeight)
	pc1.NullifyLSB()
	pc2 := NewPixelColorsFromImage(resultImg, rWidth, rHeight)

	bs := extractSecretInfoBitString(infoLength, pc1, pc2)

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
	var bb bytes.Buffer
	log.Println(bitString)
	for i := 0; i < len(bitString)/8; i++ {
		sub := bitString[i*8 : i*8+8]
		intSub, _ := strconv.ParseInt(sub, 2, 8)
		bb.WriteByte(byte(intSub))
	}

	return bb.String()
}

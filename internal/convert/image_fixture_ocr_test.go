//go:build ocr

package convert

import (
	"bytes"
	"image"
	"image/color"
	imagedraw "image/draw"
	"image/png"
	"strings"
	"testing"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func makeOCRPNG(t *testing.T, text string) []byte {
	t.Helper()

	base := image.NewRGBA(image.Rect(0, 0, 320, 80))
	imagedraw.Draw(base, base.Bounds(), &image.Uniform{C: color.White}, image.Point{}, imagedraw.Src)

	drawer := &font.Drawer{
		Dst:  base,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(20, 45),
	}
	drawer.DrawString(text)

	scaled := image.NewRGBA(image.Rect(0, 0, 1280, 320))
	xdraw.NearestNeighbor.Scale(scaled, scaled.Bounds(), base, base.Bounds(), xdraw.Src, nil)

	var pngData bytes.Buffer
	if err := png.Encode(&pngData, scaled); err != nil {
		t.Fatalf("png.Encode returned error: %v", err)
	}

	return pngData.Bytes()
}

func normalizeOCRText(text string) string {
	text = strings.ToUpper(text)
	text = strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		default:
			return ' '
		}
	}, text)

	return strings.Join(strings.Fields(text), " ")
}

package llamaimage

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type LlamaFont struct{ *truetype.Font }

func OpenFont(fontBytes []byte) (*LlamaFont, error) {
	fontStyle, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	return &LlamaFont{fontStyle}, nil
}

func (f *LlamaFont) NewFace(size float64) font.Face {
	return truetype.NewFace(f.Font, &truetype.Options{
		Size: size,
		DPI:  72,
	})
}

func (f *LlamaFont) GetWidth(fontSize float64, text string) int {
	return (&font.Drawer{
		Dst:  nil,
		Src:  nil,
		Face: f.NewFace(fontSize),
		Dot:  fixed.P(0, 0),
	}).MeasureString(text).Round()
}

func (f *LlamaFont) GetHeight(fontSize float64) int {
	return f.NewFace(fontSize).Metrics().Ascent.Ceil()
}

func (f *LlamaFont) GetTextSize(fontSize float64, text string) (width, height int) {
	width = f.GetWidth(fontSize, text)
	height = f.GetHeight(fontSize)
	return
}

func (f *LlamaFont) FitTextWidth(text string, fontSize float64, maxWidth int) (font.Face, int) {
	width := f.GetWidth(fontSize, text)
	for width >= maxWidth {
		fontSize--
		width = f.GetWidth(fontSize, text)
	}

	return f.NewFace(fontSize), width
}

func (f *LlamaFont) FitHeight(maxHeight int) (font.Face, int) {
	fontSize := float64(maxHeight)
	height := f.GetHeight(fontSize)

	for height >= maxHeight {
		fontSize--
		height = f.GetHeight(fontSize)
	}

	return f.NewFace(fontSize), height
}

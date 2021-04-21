package llamaimage

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type RawFont struct{ *truetype.Font }

func OpenFont(fontBytes []byte) (*RawFont, error) {
	fontStyle, err := truetype.Parse(fontBytes)

	if err != nil {
		return nil, err
	}
	return &RawFont{fontStyle}, nil
}

func (RFont *RawFont) NewFace(size float64) font.Face {
	return truetype.NewFace(RFont.Font, &truetype.Options{
		Size: size,
		DPI:  72,
	})
}

func (RFont *RawFont) GetTextWidth(text string, size float64) int {
	return (&font.Drawer{
		Dst:  nil,
		Src:  nil,
		Face: RFont.NewFace(size),
		Dot:  fixed.P(0, 0),
	}).MeasureString(text).Round()
}

func (RFont *RawFont) GetTextHeight(fontSize float64) int {
	return RFont.NewFace(fontSize).Metrics().Ascent.Ceil()
}

func (RFont *RawFont) GetTextSize(text string, fontSize float64) (width, height int) {
	width = RFont.GetTextWidth(text, fontSize)
	height = RFont.GetTextHeight(fontSize)
	return
}

func (RFont *RawFont) FitText(text string, fontSize float64, maxWidth int) (font.Face, int) {
	width := RFont.GetTextWidth(text, fontSize)
	for width >= maxWidth {
		fontSize--
		width = RFont.GetTextWidth(text, fontSize)
	}
	return RFont.NewFace(fontSize), width
}

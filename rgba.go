package llamaimage

import (
	"image/color"
)

type Gradient struct {
	R, G, B, A                 float64
	diffR, diffG, diffB, diffA float64
}

func GetGradientColors(of, with color.RGBA, radius float64) *Gradient {
	FR, FG, FB, FA := float64(of.R), float64(of.G), float64(of.B), float64(of.A)
	SR, SG, SB, SA := float64(with.R), float64(with.G), float64(with.B), float64(with.A)

	return &Gradient{
		FR, FG, FB, FA,
		(SR - FR) / radius,
		(SG - FG) / radius,
		(SB - FB) / radius,
		(SA - FA) / radius,
	}
}

func (d *Gradient) At(radius float64) color.RGBA {
	return color.RGBA{
		R: uint8(d.R + (d.diffR * radius)),
		G: uint8(d.G + (d.diffG * radius)),
		B: uint8(d.B + (d.diffB * radius)),
		A: uint8(d.A + (d.diffA * radius)),
	}
}

// Merges two color.RGBA color
func MergeRGBA(main, overlay color.RGBA) color.RGBA {
	oA := float64(overlay.A)

	if oA == 255 {
		return overlay
	}
	return color.RGBA{
		uint8((float64(main.R)*(255-oA) + float64(overlay.R)*oA) / 255),
		uint8((float64(main.G)*(255-oA) + float64(overlay.G)*oA) / 255),
		uint8((float64(main.B)*(255-oA) + float64(overlay.B)*oA) / 255),
		uint8(255 - (255-float64(main.A))*(255-oA)),
	}
}

// Convers hex color to color.RGBA
//
// valid format: #000000 || #000
//
// returns an error on inavlid fomat
func HexToRGBA(hexCode string) (c color.RGBA, err error) {
	c.A = 0xff

	if hexCode[0] != '#' {
		return c, ErrInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = ErrInvalidFormat
		return 0
	}

	switch len(hexCode) {
	case 7:
		c.R = hexToByte(hexCode[1])<<4 + hexToByte(hexCode[2])
		c.G = hexToByte(hexCode[3])<<4 + hexToByte(hexCode[4])
		c.B = hexToByte(hexCode[5])<<4 + hexToByte(hexCode[6])
	case 4:
		c.R = hexToByte(hexCode[1]) * 17
		c.G = hexToByte(hexCode[2]) * 17
		c.B = hexToByte(hexCode[3]) * 17
	default:
		err = ErrInvalidFormat
	}
	return
}

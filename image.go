package llamaimage

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"math"
	"os"

	_ "golang.org/x/image/webp"

	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type GradientOrientation uint8

const (
	GradientOrientationHorizontal GradientOrientation = iota
	GradientOrientationVertical
)

type ResizeMode uint8

const (
	ResizeFit ResizeMode = iota
	ResizeFill
)

var ErrInvalidFormat = errors.New("invalid format")

func NewImage(width, height int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, width, height))
}

func FillColor(img draw.Image, colorData color.RGBA) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, colorData)
		}
	}
}

func FillGradient(img draw.Image, startColor, endColor color.RGBA, orientation GradientOrientation) {
	var column, row int
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	switch orientation {
	case GradientOrientationHorizontal:
		row = height
		column = width
	case GradientOrientationVertical:
		row = width
		column = height
	}

	SR, SG, SB, SA := float64(startColor.R), float64(startColor.G), float64(startColor.B), float64(startColor.A)
	LR, LG, LB, LA := float64(endColor.R), float64(endColor.G), float64(endColor.B), float64(endColor.A)

	difference_R := (LR - SR) / float64(column)
	difference_G := (LG - SG) / float64(column)
	difference_B := (LB - SB) / float64(column)
	difference_A := (LA - SA) / float64(column)

	R, G, B, A := SR, SG, SB, SA

	for columnP := 0; columnP < column; columnP++ {
		for rowP := 0; rowP < row; rowP++ {
			switch orientation {
			case GradientOrientationHorizontal:
				img.Set(columnP, rowP, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
			case GradientOrientationVertical:
				img.Set(rowP, columnP, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
			}
		}

		R += difference_R
		G += difference_G
		B += difference_B
		A += difference_A
	}
}

func MixRGBA(first, second color.RGBA) color.RGBA {
	Alpha := 1 - (1-second.A)*(1-first.A)
	return color.RGBA{
		uint8(math.Round(float64((second.R * second.A / Alpha) + (first.R * first.A * (1 - second.A) / Alpha)))),
		uint8(math.Round(float64((second.G * second.A / Alpha) + (first.G * first.A * (1 - second.A) / Alpha)))),
		uint8(math.Round(float64((second.B * second.A / Alpha) + (first.B * first.A * (1 - second.A) / Alpha)))),
		Alpha,
	}
}

func Paste(img draw.Image, overlay image.Image, X, Y int) {
	draw.Draw(img, overlay.Bounds().Add(image.Point{X, Y}), overlay, image.Point{}, draw.Over)
}

func Write(img draw.Image, text string, textColor color.Color, fontStyle font.Face, X, Y int) {
	canvas := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: fontStyle,
		Dot:  fixed.P(X, Y+fontStyle.Metrics().Ascent.Ceil()),
	}
	canvas.DrawString(text)
}

func OpenImage(imageBytes io.Reader) (decodedImage image.Image, err error) {
	decodedImage, _, err = image.Decode(imageBytes)
	return
}

func OpenImageByPath(imagePath string) (decodedImage image.Image, err error) {
	imageBytes, err := os.Open(imagePath)
	if err != nil {
		return
	}
	defer imageBytes.Close()

	decodedImage, _, err = image.Decode(imageBytes)
	return
}

func OpenImageByBytes(imageBytes []byte) (decodedImage image.Image, err error) {
	decodedImage, _, err = image.Decode(bytes.NewReader(imageBytes))
	return
}

func OpenImageFromEFS(fileStorage fs.FS, path string) (decodedImage image.Image, err error) {
	imageBytes, err := fileStorage.Open(path)
	if err != nil {
		return nil, err
	}
	defer imageBytes.Close()
	decodedImage, _, err = image.Decode(imageBytes)
	return
}

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

func Resize(mainImage image.Image, width, height float64, mode ResizeMode) image.Image {
	var X, Y uint
	switch mode {
	case ResizeFill:
		if mainImage.Bounds().Dx() > mainImage.Bounds().Dy() {
			Y = uint(height)
		} else {
			X = uint(width)
		}
	case ResizeFit:
		imageWidth := float64(mainImage.Bounds().Dx())
		imageHeight := float64(mainImage.Bounds().Dy())
		ratio := math.Min(width/imageWidth, height/imageHeight)
		X, Y = uint(imageWidth*ratio), uint(imageHeight*ratio)
	default:
		return mainImage
	}
	return resize.Resize(X, Y, mainImage, resize.NearestNeighbor)
}

func Save(mainImage image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = png.Encode(file, mainImage); err != nil {
		return err
	}
	return nil
}

func SaveToStream(mainImage image.Image, writer io.Writer) error {
	return png.Encode(writer, mainImage)
}

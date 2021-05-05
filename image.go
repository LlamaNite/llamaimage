package llamaimage

import (
	"bytes"
	"embed"
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type GradientOrientation uint8

const (
	GradientOrientationHorizontal GradientOrientation = iota
	GradientOrientationVertical
)

type RawImage struct{ *image.RGBA }

var ErrInvalidFormat = errors.New("invalid format")

func New(width, height int) RawImage {
	generatedImage := image.NewRGBA(image.Rect(0, 0, width, height))
	return RawImage{generatedImage}
}

func (mainImage *RawImage) FillColor(colorCode color.RGBA) {
	width := mainImage.Rect.Dx()
	height := mainImage.Rect.Dy()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			mainImage.SetRGBA(x, y, colorCode)
		}
	}
}

func (mainImage *RawImage) FillGradient(startPoint, endPoint color.RGBA, style GradientOrientation) {
	var column, row int
	width, height := mainImage.Rect.Dx(), mainImage.Rect.Dy()

	switch style {
	case GradientOrientationHorizontal:
		row = height
		column = width
	case GradientOrientationVertical:
		row = width
		column = height
	}

	SR, SG, SB, SA := float64(startPoint.R), float64(startPoint.G), float64(startPoint.B), float64(startPoint.A)
	LR, LG, LB, LA := float64(endPoint.R), float64(endPoint.G), float64(endPoint.B), float64(endPoint.A)

	difference_R := (LR - SR) / float64(column)
	difference_G := (LG - SG) / float64(column)
	difference_B := (LB - SB) / float64(column)
	difference_A := (LA - SA) / float64(column)

	R, G, B, A := SR, SG, SB, SA

	for columnP := 0; columnP < column; columnP++ {
		for rowP := 0; rowP < row; rowP++ {
			switch style {
			case GradientOrientationHorizontal:
				mainImage.SetRGBA(columnP, rowP, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
			case GradientOrientationVertical:
				mainImage.SetRGBA(rowP, columnP, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
			}
		}

		R += difference_R
		G += difference_G
		B += difference_B
		A += difference_A
	}
}

func (mainImage *RawImage) Paste(overlay image.Image, X, Y int) {
	draw.Draw(mainImage, overlay.Bounds().Add(image.Point{X, Y}), overlay, image.Point{}, draw.Over)
}

func (mainImage *RawImage) Write(text string, textColor color.Color, fontStyle font.Face, X, Y int) {
	canvas := &font.Drawer{
		Dst:  mainImage,
		Src:  image.NewUniform(textColor),
		Face: fontStyle,
		Dot:  fixed.P(X, Y+fontStyle.Metrics().Ascent.Ceil()),
	}
	canvas.DrawString(text)
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

	decodedImage, err = OpenImage(imageBytes)
	return
}

func OpenImageByBytes(imageBytes []byte) (image.Image, error) {
	return OpenImage(bytes.NewReader(imageBytes))
}

func OpenImageFromEFS(fileStorage embed.FS, path string) (image.Image, error) {
	imageBytes, err := fileStorage.Open(path)
	if err != nil {
		return nil, err
	}
	return OpenImage(imageBytes)
}

func Resize(mainImage image.Image, width, height float64) image.Image {
	// return resize.Resize(uint(width), uint(height), mainImage, resize.Lanczos3)
	imageWidth := float64(mainImage.Bounds().Dx())
	imageHeight := float64(mainImage.Bounds().Dy())
	ratio := math.Min(width/imageWidth, height/imageHeight)
	return resize.Resize(uint(imageWidth*ratio), uint(imageHeight*ratio), mainImage, resize.Lanczos3)
}

func Save(mainImage image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, mainImage)
	if err != nil {
		return err
	}

	return nil
}

func SaveToStream(mainImage image.Image, writer io.Writer) error {
	return png.Encode(writer, mainImage)
}

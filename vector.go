package llamaimage

import (
	"image"
	"image/draw"

	"golang.org/x/image/vector"
)

type Vector struct {
	v *vector.Rasterizer
}

func NewVector(width, height int) *Vector {
	return &Vector{vector.NewRasterizer(width, height)}
}

func (v *Vector) Width() int {
	return v.v.Size().X
}

func (v *Vector) Height() int {
	return v.v.Size().Y
}

func (v *Vector) From(x, y int) *Vector {
	v.v.MoveTo(float32(x), float32(y))
	return v
}

func (v *Vector) To(x, y int) *Vector {
	v.v.LineTo(float32(x), float32(y))
	return v
}

func (v *Vector) Reset() *Vector {
	v.v.Reset(v.v.Size().X, v.v.Size().Y)
	return v
}

func (v *Vector) Draw(on draw.Image, with image.Image, onX, onY int) *Vector {
	v.v.ClosePath()
	v.v.Draw(on, v.v.Bounds().Add(image.Pt(onX, onY)), with, image.Point{})
	return v
}

func (v *Vector) DrawX(on draw.Image, with image.Image, onX, onY int) *Vector {
	v.v.ClosePath()
	v.v.Draw(on, v.v.Bounds().Add(image.Pt(onX, onY)), with, image.Point{})
	v.Reset()
	return v
}

func (v *Vector) DrawOp(on draw.Image, with image.Image, onX, onY int) *Vector {
	v.v.ClosePath()
	v.v.DrawOp.Draw(on, v.v.Bounds().Add(image.Pt(onX, onY)), with, image.Point{})
	return v
}

func (v *Vector) DrawOpX(on draw.Image, with image.Image, onX, onY int) *Vector {
	v.v.ClosePath()
	v.v.DrawOp.Draw(on, v.v.Bounds().Add(image.Pt(onX, onY)), with, image.Point{})
	v.Reset()
	return v
}

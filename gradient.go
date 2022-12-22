package llamaimage

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

type Point struct{ X, Y float64 }

// Gets distance between to points
//
// formula => (x - C.X) ^ 2 + (y - C.Y) ^ 2 = R^2
// returns => R
func getDistance(c, p Point) float64 {
	return math.Sqrt(math.Pow(p.X-c.X, 2) + math.Pow(p.Y-c.Y, 2))
}

// Finds the longest distance between the center point and four courners of the image
func getLongestDistance(mainImage image.Image, center Point) float64 {
	return math.Max(
		math.Max(
			getDistance(center, Point{0, 0}),
			getDistance(center, Point{float64(mainImage.Bounds().Dx()), 0}),
		),
		math.Max(
			getDistance(center, Point{0, float64(mainImage.Bounds().Dy())}),
			getDistance(center, Point{float64(mainImage.Bounds().Dx()), float64(mainImage.Bounds().Dy())}),
		),
	)
}

// Draws a Radial Gradient on the image with respect to transparent colors
// and the center point
func DrawRadialGradient(mainImage *image.RGBA, center Point, from, to color.RGBA) {
	colour := GetGradientColors(from, to, getLongestDistance(mainImage, center))

	var wg sync.WaitGroup
	var queue = make(chan int, runtime.NumCPU()*2)

	// It sets colors from top to bottom with respect to
	// it's center point at the X position update
	verticalGenerator := func() {
		for x := range queue {
			for y := 0; y < mainImage.Bounds().Dy(); y++ {
				mainImage.Set(x, y, MergeRGBA(mainImage.RGBAAt(x, y), colour.At(getDistance(center, Point{float64(x), float64(y)}))))
			}

			wg.Done()
		}
	}

	for i := 0; i < runtime.NumCPU()*2; i++ {
		go verticalGenerator()
	}

	for x := 0; x < mainImage.Bounds().Dx(); x++ {
		wg.Add(1)
		queue <- x
	}

	wg.Wait()
	close(queue)
}

// ToDo:
// func DrawLinearGradient(mainImage *image.RGBA, from, to color.RGBA) {}

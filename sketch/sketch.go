package sketch

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type SketchParams struct {
	DestWidth                int
	DestHeight               int
	StrokeRatio              float64
	StrokeReduction          float64
	StrokeJitter             int
	StrokeInversionThreshold float64
	InitialAlpha             float64
	AlphaIncrease            float64
	MinEdgeCount             int
	MaxEdgeCount             int
}

type Sketch struct {
	SketchParams      // embed for easier access
	source            image.Image
	dc                *gg.Context
	sourceWidth       int
	sourceHeight      int
	strokeSize        float64
	initialStrokeSize float64
}

func NewSketch(source image.Image, userParams SketchParams) *Sketch {
	s := &Sketch{
		SketchParams: userParams,
	}

	bounds := source.Bounds()
	s.sourceWidth, s.sourceHeight = bounds.Max.X, bounds.Max.Y
	s.initialStrokeSize = s.StrokeRatio * float64(s.DestWidth)
	s.strokeSize = s.initialStrokeSize

	canvas := gg.NewContext(s.DestWidth, s.DestHeight)
	canvas.SetColor(color.Opaque)
	for i := 0; i < 360; i += 15 {
		canvas.Push()
		canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
		canvas.DrawImage(source, rand.Int(), rand.Int())
		canvas.DrawRoundedRectangle(float64(s.DestWidth), float64(s.DestHeight), rand.Float64(), rand.Float64()*1, float64(1))
		canvas.Pop()
	}
	for j := 0; j < 100; j++ {
		for i := 0; i < 10; i++ {
			x := float64(i)*100 + 10
			y := float64(j)*100 + 50
			a1 := rand.Float64() * 2 * math.Pi
			a2 := a1 + rand.Float64()*math.Pi + math.Pi/2
			canvas.DrawArc(x, y, 40, a1, a2)
			canvas.ClosePath()
		}
	}
	canvas.FillPreserve()
	s.source = source
	s.dc = canvas
	return s
}

func (s *Sketch) Update() {

	// The core drawing logic of our sketch
	// 1. Obtain color information form the source
	rndX := rand.Float64() * float64(s.sourceWidth)
	rndY := rand.Float64() * float64(s.sourceHeight)
	r, g, b := rgb255(s.source.At(int(rndX), int(rndY)))

	// 2. Determine a destination in the output space
	destX := rndX * float64(s.DestWidth) / float64(s.sourceWidth)
	destX += float64(randRange(s.StrokeJitter))
	destY := rndY * float64(s.DestHeight) / float64(s.sourceHeight)
	destY += float64(randRange(s.StrokeJitter))
	for i := 0; i < 360; i += 15 {
		// 3. Draw a "stroke" using the desired parameters
		edges := s.MinEdgeCount + rand.Intn(s.MaxEdgeCount-s.MinEdgeCount+1)
		s.dc.SetRGBA255(r, g, b, int(s.InitialAlpha))
		s.dc.DrawRegularPolygon(edges, destX, destY, s.strokeSize, rand.Float64())
		s.dc.CubicTo(destX, destY, destX, rndY, rndX, rndY)
		s.dc.RotateAbout(gg.Radians(float64(12)), rndX/2, rndY/2)
		s.dc.Push()
		s.dc.RotateAbout(gg.Radians(float64(i)), float64(edges)/2, float64(edges)/2)
		s.dc.DrawEllipse(float64(edges)/2, float64(edges)/2, float64(edges)*7/16, float64(edges)/8)
		s.dc.Fill()

		s.dc.FillPreserve()

		if s.strokeSize <= s.StrokeInversionThreshold*s.initialStrokeSize {
			if (r+g+b)/3 < 128 {
				s.dc.SetRGBA255(255, 255, 255, int(s.InitialAlpha*2))
			} else {
				s.dc.SetRGBA255(0, 0, 0, int(s.InitialAlpha*2))
			}
		}
	}
	s.dc.Stroke()
	// 4. Update the parameter state for the next executive
	s.strokeSize -= s.StrokeReduction * s.strokeSize
	s.InitialAlpha += s.AlphaIncrease
}

func (s *Sketch) OutPut() image.Image {
	// hides the interanl implementation of this sketch
	return s.dc.Image()
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 255), int(g0 / 255), int(b0 / 255)
}

func randRange(max int) int {
	return -max + rand.Intn(2*max)
}

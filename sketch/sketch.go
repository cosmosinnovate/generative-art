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
	canvas.SetColor(color.Black)
	canvas.DrawRectangle(0, 0, float64(s.DestWidth), float64(s.DestHeight))
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

	// 3. Draw a "stroke" using the desired parameters
	edges := s.MinEdgeCount + rand.Intn(s.MaxEdgeCount-s.MinEdgeCount+1)
	s.dc.SetRGBA255(r, g, b, int(s.InitialAlpha))
	s.dc.DrawRegularPolygon(edges, math.Sin(destX), destY, s.strokeSize, rand.Float64())
	s.dc.CubicTo(math.Sin(destX), math.Atanh(destY), math.Sin(destX), math.Cos(destY), rndX, rndY)
	s.dc.DrawRoundedRectangle(destX, destY, rand.Float64(), float64(1), float64(200))
	// s.dc.RotateAbout(gg.Radians(float64(12)), rndX/2, rndY/2)
	// s.dc.DrawEllipse(math.Cos(destX), math.Sin(destY), rndX*7/16, rndY/8)
	s.dc.Stroke()
	s.dc.FillPreserve()

	if s.strokeSize <= s.StrokeInversionThreshold*s.initialStrokeSize {
		if (r+g+b)/3 < 128 {
			s.dc.SetRGBA255(255, 255, 255, int(s.InitialAlpha*2))
		} else {
			s.dc.SetRGBA255(0, 0, 0, int(s.InitialAlpha*2))
		}
	}
	
	// s.dc.Stroke()
	// 4. Update the parameter state for the next executive
	s.strokeSize -= s.StrokeReduction * s.strokeSize
	s.InitialAlpha += s.AlphaIncrease
}

// hides the interanl implementation of this sketch
func (s *Sketch) OutPut() image.Image {
	return s.dc.Image()
}

func rgb255(c color.Color) (r, g, b int) {
	r0, g0, b0, _ := c.RGBA()
	return int(r0 / 255), int(g0 / 255), int(b0 / 255)
}

func randRange(max int) int {
	return -max + rand.Intn(2*max)
}

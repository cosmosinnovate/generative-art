package main

import (
	"fmt"
	"generative-art/sketch"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
)

// var randomName = rand.Intn(99)

var (
	// sourceImgName   = "unique.jpeg"
	outputImgName   = "./out-nft/" + uuid.New().String() + "-out.png"
	totalCycleCount = 5000
)

func main() {
	rand.Seed(time.Now().UnixNano())
	img, err := randomImage(2000, 2000)
	if err != nil {
		log.Panicln(err)
	}
	destWidth := 2000
	sketch := sketch.NewSketch(img, sketch.SketchParams{
		DestWidth:                destWidth,
		DestHeight:               2000,
		StrokeRatio:              0.75,
		StrokeReduction:          0.002,
		StrokeInversionThreshold: 0.05,
		StrokeJitter:             int(0.1 * float64(destWidth)),
		InitialAlpha:             0.1,
		AlphaIncrease:            0.06,
		MinEdgeCount:             1,
		MaxEdgeCount:             4,
	})
	// the main loop
	for i := 0; i < totalCycleCount; i++ {
		sketch.Update()
	}
	saveOutput(sketch.OutPut(), outputImgName)
	// testDrawing()
}

// Test the differen functions given by gg library
func testDrawing() {
	const S = 2000
	dc := gg.NewContext(S, S)
	dc.SetRGBA(0, 0, 0, 0.1)

	for i := 0; i < 360; i += 15 {
		dc.Push()
		// dc.RotateAbout(gg.Radians(float64(i)), S/2, S/2)
		dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
		dc.Fill()
		dc.Pop()
	}
	fmt.Print(outputImgName)
	dc.SavePNG(outputImgName)
}

// Loads image from file system
// func loadImage(filePath string) (image.Image, error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("source image could not be loaded: %w", err)
// 	}
// 	defer file.Close()
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("source image could not be loaded: %w", err)
// 	}
// 	return img, nil
// }

// Load image from unsplash website. The images are retrieved randomly
func randomImage(width, height int) (image.Image, error) {
	url := fmt.Sprintf("https://source.unsplash.com/random/%dx%d", width, height)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	return img, err
}

// Save image to any filesystem
func saveOutput(img image.Image, filePath string) error {
	fmt.Println(filePath)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()
	// Encode to 'PNG' with 'DefaultCompression' level
	// then save to file
	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}

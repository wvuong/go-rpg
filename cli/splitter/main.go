package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {
	// command line arguments
	file := flag.String("file", "", "The file to split")
	x := flag.Int("x", 0, "The x coordinate to split at")
	y := flag.Int("y", 0, "The y coordinate to split at")
	width := flag.Int("width", 64, "The width of the image to split")
	height := flag.Int("height", 64, "The height of the image to split")
	frames := flag.Int("frames", 1, "The number of frames to split")
	out := flag.String("out", "frame", "Output filename prefix; frames are written as <prefix>-000.png")

	// parse command line arguments
	flag.Parse()

	if *file == "" {
		log.Fatal("-file is required")
	}

	// open image file
	fileHandle, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer fileHandle.Close()

	// read and parse image
	img, err := png.Decode(fileHandle)
	if err != nil {
		log.Fatal(err)
	}

	// check if image supports SubImage method
	si, ok := img.(SubImager)
	if !ok {
		log.Fatal("Image type does not support SubImage method")
	}

	for i := range *frames {
		// construct rect to split
		frameX := *x + (i * *width)
		rect := image.Rect(frameX, *y, frameX+*width, *y+*height)

		// a rect reaching past the sheet would silently encode a blank frame
		if !rect.In(img.Bounds()) {
			log.Fatalf("frame %d rect %v is outside the image bounds %v", i, rect, img.Bounds())
		}

		// split the image and save it under its own name — writing every frame to
		// one fixed filename meant each frame overwrote the last
		name := fmt.Sprintf("%s-%03d.png", *out, i)
		if err := writeFrame(name, si.SubImage(rect)); err != nil {
			log.Fatal(err)
		}

		log.Println("wrote", name, rect)
	}
}

func writeFrame(name string, frame image.Image) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	// close in its own scope rather than deferring inside the caller's loop,
	// where handles would pile up until main returned
	defer f.Close()

	if err := png.Encode(f, frame); err != nil {
		return err
	}

	return f.Close()
}

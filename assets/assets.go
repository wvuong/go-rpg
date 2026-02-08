package assets

import (
	"embed"
	_ "embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed *
var assets embed.FS

var Tiles_png = mustLoadImage("tiles.png")
var Sprite_png = mustLoadImage("universal-lpc-sprite_male_01_walk-3frame.png")

var Regular_ttf = mustLoadFont("mplus-1p-regular.ttf", 16)

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadFont(path string, size float64) text.Face {
	fontFile, err := assets.Open(path)
	if err != nil {
		panic(err)
	}

	s, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		panic(err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}
}

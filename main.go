package main

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wvuong/gogame/game"
	"github.com/wvuong/gogame/images"
)

var (
	tilesImage   *ebiten.Image
	spritesImage *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.Sprite_png))
	if err != nil {
		log.Fatal(err)
	}
	spritesImage = ebiten.NewImageFromImage(img)
}

func main() {
	g := game.NewGame(tilesImage, spritesImage)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

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
	tileSheet   *ebiten.Image
	spriteSheet *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tileSheet = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.Sprite_png))
	if err != nil {
		log.Fatal(err)
	}
	spriteSheet = ebiten.NewImageFromImage(img)
}

func main() {
	g := game.NewGame(tileSheet, spriteSheet)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

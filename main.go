package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	assets "github.com/wvuong/gogame/assets"
	"github.com/wvuong/gogame/game"
)

var (
	tileSheet   *ebiten.Image
	spriteSheet *ebiten.Image
)

func init() {
	tileSheet = assets.Tiles_png
	spriteSheet = assets.Sprite_png
}

func main() {
	g := game.NewGame(tileSheet, spriteSheet)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

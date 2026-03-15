package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wvuong/gogame/assets"
	"github.com/wvuong/gogame/engine"
	"github.com/wvuong/gogame/game"
)

func main() {
	assets.MustLoadAssets()

	config := engine.GameConfig{
		ScreenWidth:  1024,
		ScreenHeight: 768,
	}

	state := engine.GameState{
		WorldPosition: *engine.NewVector(100, 100),
	}

	g := game.NewGame(config, &state)

	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

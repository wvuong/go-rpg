package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wvuong/gogame/engine"
	"github.com/wvuong/gogame/game"
)

func main() {
	config := engine.GameConfig{
		ScreenWidth:  512,
		ScreenHeight: 512,
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

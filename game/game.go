package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wvuong/gogame/engine"
	"github.com/wvuong/gogame/scene"
)

type Game struct {
	state        *engine.GameState
	director     *scene.Director
	screenWidth  int
	screenHeight int
}

func NewGame(config engine.GameConfig, state *engine.GameState) Game {
	g := Game{
		state:        state,
		director:     scene.NewDirector(config, state),
		screenWidth:  config.ScreenWidth,
		screenHeight: config.ScreenHeight,
	}

	g.director.SwitchToTitle()

	return g
}

func (g *Game) Update() error {
	g.director.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.director.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	if g.screenWidth == 0 || g.screenHeight == 0 {
		return width, height
	}
	return g.screenWidth, g.screenHeight
}

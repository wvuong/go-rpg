package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wvuong/gogame/engine"
	"github.com/wvuong/gogame/scene"
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Director struct {
	scene Scene

	gameScene scene.GameScene
}

func NewDirector(config engine.GameConfig) Director {
	d := Director{
		gameScene: *scene.NewGameScene(config),
	}

	return d
}

func (d *Director) SwitchToGame() {
	d.scene = &d.gameScene
}

type Game struct {
	director     Director
	screenWidth  int
	screenHeight int
}

func NewGame(config engine.GameConfig) *Game {
	g := &Game{
		director:     NewDirector(config),
		screenWidth:  config.ScreenWidth,
		screenHeight: config.ScreenHeight,
	}

	g.director.SwitchToGame()

	return g
}

func (g *Game) Update() error {
	g.director.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.director.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	if g.screenWidth == 0 || g.screenHeight == 0 {
		return width, height
	}
	return g.screenWidth, g.screenHeight
}

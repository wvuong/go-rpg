package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wvuong/gogame/engine"
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Director struct {
	config engine.GameConfig
	state  *engine.GameState
	scene  Scene
}

func NewDirector(config engine.GameConfig, state *engine.GameState) *Director {
	d := Director{config: config, state: state}

	return &d
}

func (d *Director) Update() error {
	d.scene.Update()
	return nil
}

func (d *Director) Draw(screen *ebiten.Image) {
	d.scene.Draw(screen)
}

func (d *Director) SwitchToTitle() {
	d.scene = NewTitleScene(d.config, d.state, d)
}

func (d *Director) SwitchToTileMap() {
	d.scene = NewTileMapScene(d.config, d.state, d)
}

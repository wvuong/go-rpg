package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/wvuong/gogame/assets"
	"github.com/wvuong/gogame/engine"
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Director struct {
	config   engine.GameConfig
	state    *engine.GameState
	tileMaps map[string]*engine.TileMap

	scene Scene
}

func NewDirector(config engine.GameConfig, state *engine.GameState) *Director {
	return &Director{config: config, state: state, tileMaps: assets.TileMaps}
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
	d.scene = NewTileMapScene(d.config, d.state, d, d.tileMaps["levels/world.tmx"])
}

func (d *Director) SwitchToBattle() {
	d.scene = NewBattleScene(d.config, d.state, d)
}

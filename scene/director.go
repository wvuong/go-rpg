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
	// load tile map and assets
	layer1 := []int{
		3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, // 0
		3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, // 1
		3, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 3, // 2
		3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, // 3
		3, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 3, // 4
		3, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3, // 5
		3, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3, // 6
		3, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 3, // 7
		3, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 3, // 8
		3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, // 9
		3, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 3, // 10
		3, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3, // 11
		3, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3, // 12
		3, 3, 3, 3, 3, 3, 3, 3, 1, 1, 1, 3, // 13
	}
	//  0. 1. 2. 3. 4. 5. 6. 7. 8. 9. 10 11

	layers := [][]int{layer1}
	tileSheets := []*ebiten.Image{assets.Tiles_png}
	tileMap := engine.NewTileMap(tileSheets, layers, 12, 14, 64)

	tileMaps := make(map[string]*engine.TileMap)
	tileMaps["default"] = tileMap
	return &Director{config: config, state: state, tileMaps: tileMaps}
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
	d.scene = NewTileMapScene(d.config, d.state, d, d.tileMaps["default"])
}

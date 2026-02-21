package scene

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wvuong/gogame/assets"
	"github.com/wvuong/gogame/engine"
	"golang.org/x/image/colornames"
)

const (
	cols     = 12
	rows     = 14
	tileSize = 64
)

type TileMapScene struct {
	config      engine.GameConfig
	state       *engine.GameState
	director    *Director
	spriteSheet *ebiten.Image
	tileMap     *engine.TileMap
	camera      *engine.Camera
	player      *engine.Player
	debug       *engine.Debug
}

func NewTileMapScene(
	config engine.GameConfig,
	state *engine.GameState,
	director *Director) *TileMapScene {

	tileSheet := assets.Tiles_png
	spriteSheet := assets.Sprite_png

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
	tileMap := engine.NewMap(tileSheet, [][]int{layer1}, cols, rows, tileSize)

	// create sprite indexes for each direction
	frameInterval := 100 * time.Millisecond
	facingDownSpriteIndex := engine.NewSpriteIndex(spriteSheet, []image.Rectangle{
		image.Rect(0, 128, 48, 192),
		image.Rect(48, 128, 96, 192),
		image.Rect(96, 128, 144, 192),
	}, frameInterval)

	facingUpSpriteIndex := engine.NewSpriteIndex(spriteSheet, []image.Rectangle{
		image.Rect(0, 0, 48, 64),
		image.Rect(48, 0, 96, 64),
		image.Rect(96, 0, 144, 64),
	}, frameInterval)

	facingLeftSpriteIndex := engine.NewSpriteIndex(spriteSheet, []image.Rectangle{
		image.Rect(0, 192, 48, 256),
		image.Rect(48, 192, 96, 256),
		image.Rect(96, 192, 144, 256),
	}, frameInterval)

	facingRightSpriteIndex := engine.NewSpriteIndex(spriteSheet, []image.Rectangle{
		image.Rect(0, 64, 48, 128),
		image.Rect(48, 64, 96, 128),
		image.Rect(96, 64, 144, 128),
	}, frameInterval)

	directionalSpriteIndex := &engine.DirectionalSpriteIndex{
		Up:    facingUpSpriteIndex,
		Down:  facingDownSpriteIndex,
		Left:  facingLeftSpriteIndex,
		Right: facingRightSpriteIndex,
	}

	// debug state storage
	debug := &engine.Debug{Enabled: false}

	// create player at position 100, 100
	player := engine.NewPlayer(tileMap, directionalSpriteIndex, &state.WorldPosition, 2)

	// create camera and center on player
	camera := engine.NewCamera(config, tileMap, debug)
	camera.CenterOn(player.Sprite)

	g := &TileMapScene{
		config:   config,
		state:    state,
		director: director,
		tileMap:  tileMap,
		camera:   camera,
		player:   player,
		debug:    debug,
	}

	return g
}

func (g *TileMapScene) Update() {
	// go back to title screen if escape is pressed
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.director.SwitchToTitle()
	}

	// toggle debug mode
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.debug.Enabled = !g.debug.Enabled
	}

	// update player
	g.player.Update()

	// update camera position
	g.camera.Update()
}

func (g *TileMapScene) Draw(screen *ebiten.Image) {
	// draw the tilemap with the camera's viewport
	g.camera.Draw(screen)

	// draw player 48x64 sprite
	// center sprite on its position
	op := &ebiten.DrawImageOptions{}
	centeredX := g.player.Sprite.ScreenPosition.X - float64(g.player.Sprite.Dx)/2
	centeredY := g.player.Sprite.ScreenPosition.Y - float64(g.player.Sprite.Dy)/2
	op.GeoM.Translate(centeredX, centeredY)
	screen.DrawImage(g.player.Sprite.Image, op)

	if g.debug.Enabled {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.0f", g.player.Left), int(centeredX), int(centeredY))

		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("Player: (%.0f, %.0f) (%.0f, %.0f)",
				g.player.Sprite.Position.X, g.player.Sprite.Position.Y,
				g.player.Sprite.ScreenPosition.X, g.player.Sprite.ScreenPosition.Y), 0, 0)
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("(L%.0f, R%.0f, T%.0f, B%.0f)", g.player.Left, g.player.Right, g.player.Top, g.player.Bottom), 0, g.config.ScreenHeight-20)

		vector.FillCircle(screen, float32(g.player.Sprite.ScreenPosition.X), float32(g.player.Sprite.ScreenPosition.Y), 2, colornames.Red, false)

		vector.StrokeRect(screen, float32(centeredX), float32(centeredY),
			float32(g.player.Sprite.Dx), float32(g.player.Sprite.Dy), 2, color.Black, false)
	}
}

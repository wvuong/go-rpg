package scene

import (
	"fmt"
	"image"
	"image/color"
	"math"
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

type GameScene struct {
	config      engine.GameConfig
	tileSheet   *ebiten.Image
	spriteSheet *ebiten.Image
	tileMap     *engine.TileMap
	camera      *engine.Camera
	player      *engine.Player
	debug       *engine.Debug
}

func NewGameScene(config engine.GameConfig) *GameScene {
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
	tileMap := engine.NewMap([][]int{layer1}, cols, rows, tileSize)

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

	// create player at position 100, 100
	player := engine.NewPlayer(tileMap, directionalSpriteIndex, 100, 100, 2)

	// create camera and center on player
	camera := engine.NewCamera(config.ScreenWidth, config.ScreenHeight, cols, rows, tileSize)
	camera.CenterOn(player.Sprite)

	g := &GameScene{
		config:    config,
		tileSheet: tileSheet,
		tileMap:   tileMap,
		camera:    camera,
		player:    player,
		debug:     &engine.Debug{Enabled: false},
	}

	return g
}

func (g *GameScene) Update() {
	// toggle debug mode
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.debug.Enabled = !g.debug.Enabled
	}

	// update player
	g.player.Update()

	// update camera position
	g.camera.Update()
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	// use camera position to determine visible tiles
	startCol := int(math.Floor(float64(g.camera.Position.X / tileSize)))
	endCol := startCol + int(math.Floor(float64(g.camera.Width/tileSize)))
	startRow := int(math.Floor(g.camera.Position.Y / tileSize))
	endRow := startRow + int(g.camera.Height/tileSize)
	offsetX := -g.camera.Position.X + float64(startCol*tileSize)
	offsetY := -g.camera.Position.Y + float64(startRow*tileSize)

	// draw layers
	for l := range g.tileMap.Layers {
		for c := startCol; c <= endCol; c++ {
			// x is the screen position of the tile
			x := float64((c-startCol)*tileSize) + offsetX

			for r := startRow; r <= endRow; r++ {
				// y is the screen position of the tile
				y := float64((r-startRow)*tileSize) + offsetY

				// this is the raw tile index from the map data
				tileId := g.tileMap.GetTile(l, c, r)
				// 0 => empty tile
				if tileId != 0 {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(x, y)

					tileIdx := tileId - 1 // subtract 1 to get the tile array index

					// create a rect of the tile in the tilesheet
					sx := tileIdx * tileSize
					rect := image.Rect(sx, 0, sx+tileSize, tileSize)

					// if the tile is solid, tint it red
					if g.debug.Enabled && tileId > 2 {
						op.ColorScale.ScaleWithColor(colornames.Red)
					}

					// draw tile
					screen.DrawImage(g.tileSheet.SubImage(rect).(*ebiten.Image), op)

					if g.debug.Enabled {
						ebitenutil.DebugPrintAt(screen, fmt.Sprintf("(%d)", tileId),
							int(x+tileSize/2), int(y+tileSize/2))
					}
				}

				// draw horizontal grid line
				if g.debug.Enabled {
					vector.StrokeLine(screen, 0, float32(y), float32(g.config.ScreenWidth), float32(y), 2, color.Black, false)
					ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %.0f", r, y), 10, int(y))
				}
			}

			// draw vertical grid line
			if g.debug.Enabled {
				vector.StrokeLine(screen, float32(x), 0, float32(x), float32(g.config.ScreenHeight), 2, color.Black, false)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %.0f", c, x), int(x), 10)
			}
		}
	}

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

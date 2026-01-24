package game

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

const (
	cols         = 12
	rows         = 14
	tileSize     = 64
	screenWidth  = 8 * tileSize
	screenHeight = 8 * tileSize
	mapWidth     = cols * tileSize
	mapHeight    = rows * tileSize
)

type Game struct {
	tileSheet   *ebiten.Image
	spriteSheet *ebiten.Image
	tileMap     *TileMap
	camera      *Camera
	player      *Player
	debug       *Debug
}

func NewGame(tileSheet *ebiten.Image, spriteSheet *ebiten.Image) *Game {
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
	tileMap := NewMap([][]int{layer1}, cols, rows, tileSize)

	rect := image.Rect(48, 128, 96, 192)
	sprite := spriteSheet.SubImage(rect).(*ebiten.Image)

	camera := NewCamera(screenWidth, screenHeight, cols, rows, tileSize)

	player := &Player{
		Sprite: NewSprite(sprite, 100, 100),
	}

	camera.CenterOn(player.Sprite)

	g := &Game{
		tileSheet: tileSheet,
		tileMap:   tileMap,
		camera:    camera,
		player:    player,
		debug:     &Debug{Enabled: false},
	}

	return g
}

func (g *Game) Update() error {
	// move player with arrow keys
	var dirX, dirY float64
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dirY = -4
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dirY = 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dirX = -4
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dirX = 4
	}

	// toggle debug mode
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.debug.Enabled = !g.debug.Enabled
	}

	g.player.Sprite.Position.X += dirX
	g.player.Sprite.Position.Y += dirY

	// check for collisions
	// calculate the player's bounding box
	left := g.player.Sprite.Position.X - float64(g.player.Sprite.Dx)/2
	right := g.player.Sprite.Position.X + float64(g.player.Sprite.Dx)/2 - 1
	top := g.player.Sprite.Position.Y - float64(g.player.Sprite.Dy)/2
	bottom := g.player.Sprite.Position.Y + float64(g.player.Sprite.Dy)/2 - 1
	g.player.Left = left
	g.player.Right = right
	g.player.Top = top
	g.player.Bottom = bottom

	collision := g.tileMap.isSolidTileAtXY(left, top) ||
		g.tileMap.isSolidTileAtXY(right, top) ||
		g.tileMap.isSolidTileAtXY(right, bottom) ||
		g.tileMap.isSolidTileAtXY(left, bottom)

	if collision {
		if dirY > 0 {
			// moving down
			row := g.tileMap.getRow(bottom)
			// align player to top edge of tile
			g.player.Sprite.Position.Y = -float64(g.player.Sprite.Dy)/2 + g.tileMap.getY(row)
		} else if dirY < 0 {
			// moving up
			row := g.tileMap.getRow(top)
			// align player to bottom edge of tile
			g.player.Sprite.Position.Y = float64(g.player.Sprite.Dy)/2 + g.tileMap.getY(row+1)
		} else if dirX > 0 {
			// moving right
			col := g.tileMap.getCol(right)
			// align player to left edge of tile
			g.player.Sprite.Position.X = -float64(g.player.Sprite.Dx)/2 + g.tileMap.getX(col)
		} else if dirX < 0 {
			// moving left
			col := g.tileMap.getCol(left)
			// align player to right edge of tile
			g.player.Sprite.Position.X = float64(g.player.Sprite.Dx)/2 + g.tileMap.getX(col+1)
		}
	}

	// clamp player position to map bounds
	x := math.Max(0, math.Min(g.player.Sprite.Position.X, mapWidth))
	y := math.Max(0, math.Min(g.player.Sprite.Position.Y, mapHeight))
	g.player.Sprite.Position.X = x
	g.player.Sprite.Position.Y = y

	g.camera.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
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
					vector.StrokeLine(screen, 0, float32(y), float32(screenWidth), float32(y), 2, color.Black, false)
					ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %.0f", r, y), 10, int(y))
				}
			}

			// draw vertical grid line
			if g.debug.Enabled {
				vector.StrokeLine(screen, float32(x), 0, float32(x), float32(screenHeight), 2, color.Black, false)
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
			fmt.Sprintf("(L%.0f, R%.0f, T%.0f, B%.0f)", g.player.Left, g.player.Right, g.player.Top, g.player.Bottom), 0, screenHeight-20)

		vector.FillCircle(screen, float32(g.player.Sprite.ScreenPosition.X), float32(g.player.Sprite.ScreenPosition.Y), 2, colornames.Red, false)

		vector.StrokeRect(screen, float32(centeredX), float32(centeredY),
			float32(g.player.Sprite.Dx), float32(g.player.Sprite.Dy), 2, color.Black, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

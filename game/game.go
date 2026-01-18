package game

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	}

	return g
}

func (g *Game) Update() error {
	// move player with arrow keys
	var dirX, dirY float64
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		dirY = -16
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		dirY = 16
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		dirX = -16
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		dirX = 16
	}

	g.player.Sprite.Position.X += dirX
	g.player.Sprite.Position.Y += dirY

	// check for collisions
	left := int(g.player.Sprite.Position.X - float64(g.player.Sprite.Dx)/2)
	right := int(g.player.Sprite.Position.X + float64(g.player.Sprite.Dx)/2 - 1)
	top := int(g.player.Sprite.Position.Y - float64(g.player.Sprite.Dy)/2)
	bottom := int(g.player.Sprite.Position.Y + float64(g.player.Sprite.Dy)/2 - 1)
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
			row := g.tileMap.getRow(bottom)
			g.player.Sprite.Position.Y = float64(g.player.Sprite.Dy/2) + g.tileMap.getY(row)
		} else if dirY < 0 {
			row := g.tileMap.getRow(top)
			g.player.Sprite.Position.Y = float64(g.player.Sprite.Dy/2) + g.tileMap.getY(row+1)
		} else if dirX > 0 {
			col := g.tileMap.getCol(right)
			g.player.Sprite.Position.X = float64(g.player.Sprite.Dx/2) + g.tileMap.getX(col)
		} else if dirX < 0 {
			col := g.tileMap.getCol(left)
			g.player.Sprite.Position.X = float64(g.player.Sprite.Dx/2) + g.tileMap.getX(col+1)
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
			for r := startRow; r <= endRow; r++ {
				tile := g.tileMap.GetTile(l, c, r)
				if tile != 0 {
					// 0 => empty tile
					op := &ebiten.DrawImageOptions{}
					x := float64((c-startCol)*tileSize) + offsetX
					y := float64((r-startRow)*tileSize) + offsetY
					op.GeoM.Translate(x, y)

					tileIdx := tile - 1
					sx := tileIdx * tileSize
					rect := image.Rect(sx, 0, sx+tileSize, tileSize)
					screen.DrawImage(g.tileSheet.SubImage(rect).(*ebiten.Image), op)
				}
			}
		}
	}

	// draw player 48x64 sprite
	// center sprite on its position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.player.Sprite.ScreenPosition.X-float64(g.player.Sprite.Dx)/2,
		g.player.Sprite.ScreenPosition.Y-float64(g.player.Sprite.Dy)/2)
	screen.DrawImage(g.player.Sprite.Image, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player: (%.2f, %.2f) (L%d, R%d, T%d, B%d)\nCamera: (%.2f, %.2f)",
		g.player.Sprite.Position.X, g.player.Sprite.Position.Y,
		g.player.Left, g.player.Right, g.player.Top, g.player.Bottom,
		g.camera.Position.X, g.camera.Position.Y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

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
)

type Game struct {
	tilesImage   *ebiten.Image
	spritesImage *ebiten.Image
	tileMap      *TileMap
	//player       *Sprite
	camera *Camera
}

func NewGame(tilesImage *ebiten.Image, spritesImage *ebiten.Image) *Game {
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
	tileMap := NewMap([][]int{layer1}, cols, rows)

	//rect := image.Rect(48, 128, 96, 192)
	//sprite := spritesImage.SubImage(rect).(*ebiten.Image)

	camera := NewCamera(screenWidth, screenHeight, cols, rows, tileSize)

	g := &Game{
		tilesImage: tilesImage,
		tileMap:    tileMap,
		//player:     &Sprite{image: sprite, x: 100, y: 100},
		camera: camera,
	}

	return g
}

func (g *Game) Update() error {
	speed := 10.0

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		//g.player.y -= speed
		g.camera.Move(0, -speed)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		//g.player.y += speed
		g.camera.Move(0, speed)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		//g.player.x -= speed
		g.camera.Move(-speed, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		//g.player.x += speed
		g.camera.Move(speed, 0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
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
					screen.DrawImage(g.tilesImage.SubImage(rect).(*ebiten.Image), op)
				}
			}
		}
	}

	// draw 48x64 sprite
	/*
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.player.x), float64(g.player.y))
		screen.DrawImage(g.player.image, op)
	*/
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f, (%d, %d) -> (%d, %d)",
		ebiten.ActualTPS(), startCol, startRow, endCol, endRow))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

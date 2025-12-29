package main

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wvuong/gogame/images"
)

const (
	cols         = 8
	rows         = 8
	tileSize     = 64
	screenWidth  = cols * tileSize
	screenHeight = rows * tileSize
)

var (
	tilesImage  *ebiten.Image
	spriteImage *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)

	img, _, err = image.Decode(bytes.NewReader(images.Sprite_png))
	if err != nil {
		log.Fatal(err)
	}
	spriteImage = ebiten.NewImageFromImage(img)
}

type Sprite struct {
	image *ebiten.Image
	x     int
	y     int
}

type Game struct {
	layers [][]int
	player *Sprite
}

func (g *Game) Update() error {
	speed := 5

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.player.y -= speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.player.y += speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.player.x -= speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.player.x += speed
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw layers
	for _, l := range g.layers {
		for c := range cols {
			for r := range rows {
				tile := GetTile(l, c, r)
				if tile != 0 {
					// 0 => empty tile
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(c*tileSize), float64(r*tileSize))

					tileIdx := tile - 1
					sx := tileIdx * tileSize
					rect := image.Rect(sx, 0, sx+tileSize, tileSize)
					screen.DrawImage(tilesImage.SubImage(rect).(*ebiten.Image), op)
				}
			}
		}
	}

	// draw 48x64 sprite
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.player.x), float64(g.player.y))
	screen.DrawImage(g.player.image, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func GetTile(layer []int, col int, row int) int {
	return layer[row*cols+col]
}

func main() {
	layer1 := []int{
		1, 3, 3, 3, 1, 1, 3, 1, // 0
		1, 1, 1, 1, 1, 1, 1, 1, // 1
		1, 1, 1, 1, 1, 2, 1, 1, // 2
		1, 1, 1, 1, 1, 1, 1, 1, // 3
		1, 1, 1, 2, 1, 1, 1, 1, // 4
		1, 1, 1, 1, 2, 1, 1, 1, // 5
		1, 1, 1, 1, 2, 1, 1, 1, // 6
		1, 1, 1, 0, 0, 1, 1, 1, // 7
	}
	//  0. 1. 2. 3. 4. 5. 6. 7.

	layer2 := []int{
		0, 0, 0, 0, 0, 0, 0, 0, // 0
		0, 0, 0, 0, 0, 0, 0, 0, // 1
		0, 0, 0, 0, 0, 0, 0, 0, // 2
		0, 0, 0, 0, 0, 0, 0, 0, // 3
		5, 5, 0, 0, 0, 0, 4, 0, // 4
		5, 5, 0, 0, 0, 0, 3, 0, // 5
		0, 0, 0, 0, 0, 0, 0, 0, // 6
		0, 0, 0, 0, 0, 0, 0, 0, // 7
	}
	//  0. 1. 2. 3. 4. 5. 6. 7.

	rect := image.Rect(48, 128, 96, 192)
	sprite := spriteImage.SubImage(rect).(*ebiten.Image)

	g := &Game{
		layers: [][]int{layer1, layer2},
		player: &Sprite{image: sprite, x: 100, y: 100},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tiles (Ebitengine Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

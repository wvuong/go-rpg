package engine

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Camera struct {
	config      GameConfig
	tileMap     *TileMap
	debug       *Debug
	position    Vector
	maxPosition Vector
	target      *Sprite
}

func NewCamera(
	config GameConfig,
	tileMap *TileMap,
	debug *Debug) *Camera {

	return &Camera{
		config:      config,
		tileMap:     tileMap,
		debug:       debug,
		position:    Vector{X: 0, Y: 0},
		maxPosition: Vector{X: float64(tileMap.cols*tileMap.tileSize - config.ScreenWidth), Y: float64(tileMap.rows*tileMap.tileSize - config.ScreenHeight)},
	}
}

func (c *Camera) Move(dirx, diry float64) {
	c.position.X += dirx
	c.position.Y += diry

	// clamp values
	c.position.X = math.Max(0, math.Min(c.position.X, c.maxPosition.X))
	c.position.Y = math.Max(0, math.Min(c.position.Y, c.maxPosition.Y))
}

func (c *Camera) CenterOn(sprite *Sprite) {
	c.target = sprite
	sprite.ScreenPosition.X = 0
	sprite.ScreenPosition.Y = 0
}

func (c *Camera) Update() {
	if c.target != nil {
		width := float64(c.config.ScreenWidth)
		height := float64(c.config.ScreenHeight)

		// assume followed sprite should be placed at the center of the screen
		// whenever possible
		c.target.ScreenPosition.X = width / 2
		c.target.ScreenPosition.Y = height / 2

		// make the camera follow the sprite
		c.position.X = c.target.Position.X - width/2
		c.position.Y = c.target.Position.Y - height/2

		// clamp values
		c.position.X = math.Max(0, math.Min(c.position.X, c.maxPosition.X))
		c.position.Y = math.Max(0, math.Min(c.position.Y, c.maxPosition.Y))

		// in map corners, the sprite cannot be placed in the center of the screen
		// and we have to change its screen coordinates

		// left and right sides
		if c.target.Position.X < width/2 ||
			c.target.Position.X > c.maxPosition.X+width/2 {
			c.target.ScreenPosition.X = c.target.Position.X - c.position.X
		}
		// top and bottom sides
		if c.target.Position.Y < height/2 ||
			c.target.Position.Y > c.maxPosition.Y+height/2 {
			c.target.ScreenPosition.Y = c.target.Position.Y - c.position.Y
		}
	}

}

func (c *Camera) Draw(screen *ebiten.Image) {
	tileSize := c.tileMap.tileSize
	tileSizeFloat := float64(c.tileMap.tileSize)
	width := c.config.ScreenWidth
	height := c.config.ScreenHeight

	// use camera position to determine visible tiles
	startCol := int(math.Floor(c.position.X / tileSizeFloat))
	endCol := startCol + int(math.Floor(float64(width/tileSize)))
	startRow := int(math.Floor(c.position.Y / tileSizeFloat))
	endRow := startRow + int(height/tileSize)
	offsetX := -c.position.X + float64(startCol*tileSize)
	offsetY := -c.position.Y + float64(startRow*tileSize)

	// draw layers
	for layer := range c.tileMap.layers {
		for col := startCol; col <= endCol; col++ {
			// x is the screen position of the tile
			x := float64((col-startCol)*tileSize) + offsetX

			for row := startRow; row <= endRow; row++ {
				// y is the screen position of the tile
				y := float64((row-startRow)*tileSize) + offsetY

				// draw tile to screen
				c.tileMap.drawTile(screen, layer, col, row, x, y, c.debug)

				// draw horizontal grid line
				if c.debug.Enabled {
					vector.StrokeLine(screen, 0, float32(y), float32(c.config.ScreenWidth), float32(y), 2, color.Black, false)
					ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %.0f", row, y), 10, int(y))
				}
			}

			// draw vertical grid line
			if c.debug.Enabled {
				vector.StrokeLine(screen, float32(x), 0, float32(x), float32(c.config.ScreenHeight), 2, color.Black, false)
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %.0f", col, x), int(x), 10)
			}
		}
	}
}

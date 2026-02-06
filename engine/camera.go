package engine

import (
	"math"
)

type Camera struct {
	Position    Vector
	Width       int
	Height      int
	maxPosition Vector
	target      *Sprite
}

func NewCamera(width int, height int, numberColumns int, numberRows int, tileSize int) *Camera {
	return &Camera{
		Position:    Vector{X: 0, Y: 0},
		Width:       width,
		Height:      height,
		maxPosition: Vector{X: float64(numberColumns*tileSize - width), Y: float64(numberRows*tileSize - height)},
	}
}

func (c *Camera) Move(dirx, diry float64) {
	c.Position.X += dirx
	c.Position.Y += diry

	// clamp values
	c.Position.X = math.Max(0, math.Min(c.Position.X, c.maxPosition.X))
	c.Position.Y = math.Max(0, math.Min(c.Position.Y, c.maxPosition.Y))
}

func (c *Camera) CenterOn(sprite *Sprite) {
	c.target = sprite
	sprite.ScreenPosition.X = 0
	sprite.ScreenPosition.Y = 0
}

func (c *Camera) Update() {
	if c.target != nil {
		// assume followed sprite should be placed at the center of the screen
		// whenever possible
		c.target.ScreenPosition.X = float64(c.Width) / 2
		c.target.ScreenPosition.Y = float64(c.Height) / 2

		// make the camera follow the sprite
		c.Position.X = c.target.Position.X - float64(c.Width)/2
		c.Position.Y = c.target.Position.Y - float64(c.Height)/2

		// clamp values
		c.Position.X = math.Max(0, math.Min(c.Position.X, c.maxPosition.X))
		c.Position.Y = math.Max(0, math.Min(c.Position.Y, c.maxPosition.Y))

		// in map corners, the sprite cannot be placed in the center of the screen
		// and we have to change its screen coordinates

		// left and right sides
		if c.target.Position.X < float64(c.Width)/2 ||
			c.target.Position.X > c.maxPosition.X+float64(c.Width)/2 {
			c.target.ScreenPosition.X = c.target.Position.X - c.Position.X
		}
		// top and bottom sides
		if c.target.Position.Y < float64(c.Height)/2 ||
			c.target.Position.Y > c.maxPosition.Y+float64(c.Height)/2 {
			c.target.ScreenPosition.Y = c.target.Position.Y - c.Position.Y
		}
	}

}

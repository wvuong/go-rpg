package engine

import (
	"math"
)

type Camera struct {
	Position    Vector
	Width       int
	Height      int
	MaxPosition Vector
	Target      *Sprite
}

func NewCamera(width int, height int, numberColumns int, numberRows int, tileSize int) *Camera {
	return &Camera{
		Position:    Vector{X: 0, Y: 0},
		Width:       width,
		Height:      height,
		MaxPosition: Vector{X: float64(numberColumns*tileSize - width), Y: float64(numberRows*tileSize - height)},
	}
}

func (c *Camera) Move(dirx, diry float64) {
	c.Position.X += dirx
	c.Position.Y += diry

	// clamp values
	c.Position.X = math.Max(0, math.Min(c.Position.X, c.MaxPosition.X))
	c.Position.Y = math.Max(0, math.Min(c.Position.Y, c.MaxPosition.Y))
}

func (c *Camera) CenterOn(sprite *Sprite) {
	c.Target = sprite
	sprite.ScreenPosition.X = 0
	sprite.ScreenPosition.Y = 0
}

func (c *Camera) Update() {
	if c.Target != nil {
		// assume followed sprite should be placed at the center of the screen
		// whenever possible
		c.Target.ScreenPosition.X = float64(c.Width) / 2
		c.Target.ScreenPosition.Y = float64(c.Height) / 2

		// make the camera follow the sprite
		c.Position.X = c.Target.Position.X - float64(c.Width)/2
		c.Position.Y = c.Target.Position.Y - float64(c.Height)/2

		// clamp values
		c.Position.X = math.Max(0, math.Min(c.Position.X, c.MaxPosition.X))
		c.Position.Y = math.Max(0, math.Min(c.Position.Y, c.MaxPosition.Y))

		// in map corners, the sprite cannot be placed in the center of the screen
		// and we have to change its screen coordinates

		// left and right sides
		if c.Target.Position.X < float64(c.Width)/2 ||
			c.Target.Position.X > c.MaxPosition.X+float64(c.Width)/2 {
			c.Target.ScreenPosition.X = c.Target.Position.X - c.Position.X
		}
		// top and bottom sides
		if c.Target.Position.Y < float64(c.Height)/2 ||
			c.Target.Position.Y > c.MaxPosition.Y+float64(c.Height)/2 {
			c.Target.ScreenPosition.Y = c.Target.Position.Y - c.Position.Y
		}
	}

}

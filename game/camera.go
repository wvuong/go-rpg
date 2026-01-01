package game

import (
	"math"
)

type Camera struct {
	Position    Vector
	Width       int
	Height      int
	MaxPosition Vector
	Speed       float32
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

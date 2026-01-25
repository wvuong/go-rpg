package engine

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Image          *ebiten.Image
	Position       *Vector
	ScreenPosition *Vector
	Dx             int
	Dy             int
}

func NewSprite(image *ebiten.Image, x, y float64) *Sprite {
	return &Sprite{
		Image:    image,
		Position: &Vector{X: x, Y: y},
		ScreenPosition: &Vector{
			X: 0,
			Y: 0,
		},
		Dx: image.Bounds().Dx(),
		Dy: image.Bounds().Dy(),
	}
}

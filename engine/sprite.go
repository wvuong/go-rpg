package engine

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Image          *ebiten.Image
	Position       *Vector
	ScreenPosition *Vector
	Dx             int
	Dy             int
	Layer          int
}

func NewSprite(image *ebiten.Image, position *Vector) *Sprite {
	return &Sprite{
		Image:    image,
		Position: position,
		ScreenPosition: &Vector{
			X: 0,
			Y: 0,
		},
		Dx:    image.Bounds().Dx(),
		Dy:    image.Bounds().Dy(),
		Layer: 0,
	}
}

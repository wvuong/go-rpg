package engine

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	Image          *ebiten.Image
	ScreenPosition *Vector
	Dx             int
	Dy             int
	Layer          int
}

func NewSprite(image *ebiten.Image) *Sprite {
	return &Sprite{
		Image:          image,
		ScreenPosition: NewVector(0, 0),
		Dx:             image.Bounds().Dx(),
		Dy:             image.Bounds().Dy(),
		Layer:          0,
	}
}

package engine

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Sprite struct {
	Image          *ebiten.Image
	ScreenPosition *Vector
	Dx             int
	Dy             int
	HorizontalFlip bool

	// degrees, 0 is facing right
	OriginalRotation float64

	// degrees, 0 is facing right
	Rotation float64
}

func NewSprite(image *ebiten.Image) *Sprite {
	return &Sprite{
		Image:            image,
		ScreenPosition:   NewVector(0, 0),
		Dx:               image.Bounds().Dx(),
		Dy:               image.Bounds().Dy(),
		HorizontalFlip:   false,
		OriginalRotation: 0,
		Rotation:         0,
	}
}

func NewSprite2(image *ebiten.Image, horizontalFlip bool, rotation float64) *Sprite {
	return &Sprite{
		Image:            image,
		ScreenPosition:   NewVector(0, 0),
		Dx:               image.Bounds().Dx(),
		Dy:               image.Bounds().Dy(),
		HorizontalFlip:   horizontalFlip,
		OriginalRotation: rotation,
		Rotation:         0,
	}
}

func (s *Sprite) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions, debug *Debug) {
	// do horizontal flipping
	if s.HorizontalFlip {
		// if flipped, we need to adjust the position after flipping to keep it in the same place
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(s.Dx), 0)
	}

	halfW, halfH := float64(s.Dx)/2, float64(s.Dy)/2

	// do rotation
	radians := s.Rotation * 2 * math.Pi / 360
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(radians)
	op.GeoM.Translate(halfW, halfH)

	// translate into position
	centeredX := s.ScreenPosition.X - halfW
	centeredY := s.ScreenPosition.Y - halfH
	op.GeoM.Translate(centeredX, centeredY)

	// draw onto screen
	screen.DrawImage(s.Image, op)

	// draw debug info
	if debug.Enabled {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("(%.0f, %.0f) c(%.0f, %.0f) [%d, %d]", s.ScreenPosition.X, s.ScreenPosition.Y, centeredX, centeredY, s.Dx, s.Dy), int(centeredX), int(centeredY))

		vector.StrokeRect(screen, float32(centeredX), float32(centeredY),
			float32(s.Dx), float32(s.Dy), 2, color.Black, false)
	}
}

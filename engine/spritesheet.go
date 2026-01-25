package engine

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteIndex struct {
	images []*ebiten.Image
	idx    int
	timer  *Timer
}

func NewSpriteIndex(image *ebiten.Image, rects []image.Rectangle) *SpriteIndex {
	slice := make([]*ebiten.Image, len(rects))
	for i, rect := range rects {
		// process each rect
		slice[i] = image.SubImage(rect).(*ebiten.Image)
	}

	return &SpriteIndex{
		images: slice,
		idx:    0,
		timer:  NewTimer(100 * time.Millisecond),
	}
}

func (si *SpriteIndex) NextFrame() *ebiten.Image {
	next := si.images[si.idx]
	si.timer.Update()

	if si.timer.IsReady() {
		si.idx = (si.idx + 1) % len(si.images)
		si.timer.Reset()
	}

	return next
}

type DirectionalSpriteIndex struct {
	Up    *SpriteIndex
	Down  *SpriteIndex
	Left  *SpriteIndex
	Right *SpriteIndex
}

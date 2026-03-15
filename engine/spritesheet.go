package engine

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteIndex struct {
	images    []*ebiten.Image
	frames    int
	intervals []time.Duration
	idx       int
	timer     *Timer
	loop      bool
	duration  time.Duration
}

func NewSpriteIndex(spriteSheet *ebiten.Image, rects []image.Rectangle, loop bool, intervals []time.Duration) *SpriteIndex {
	slice := make([]*ebiten.Image, len(rects))
	for i, rect := range rects {
		// process each rect
		slice[i] = spriteSheet.SubImage(rect).(*ebiten.Image)
	}

	duration := 0 * time.Second
	for i := range intervals {
		duration += intervals[i]
	}

	return &SpriteIndex{
		images:    slice,
		frames:    len(rects),
		intervals: intervals,
		idx:       0,
		timer:     NewTimer(intervals[0]),
		loop:      loop,
		duration:  duration,
	}
}

func NewHorizontalSpriteIndex(spriteSheet *ebiten.Image, width int, height int, frames int, loop bool, intervals []time.Duration) *SpriteIndex {
	slice := make([]*ebiten.Image, frames)
	for i := range frames {
		rect := image.Rect(width*i, 0, width*(i+1), height)
		slice[i] = spriteSheet.SubImage(rect).(*ebiten.Image)
	}

	duration := 0 * time.Second
	for i := range intervals {
		duration += intervals[i]
	}

	return &SpriteIndex{
		images:    slice,
		frames:    frames,
		intervals: intervals,
		idx:       0,
		timer:     NewTimer(intervals[0]),
		loop:      loop,
		duration:  duration,
	}
}

func (si *SpriteIndex) Reset() {
	si.idx = 0
}

func (si *SpriteIndex) NextFrame() *ebiten.Image {
	// get the current frame
	next := si.images[si.idx]
	si.timer.Update()

	// if the timer is ready to advance to the next frame
	if si.timer.IsReady() {
		if si.loop {
			// move to the next frame
			si.idx = (si.idx + 1) % si.frames

		} else {
			si.idx = si.idx + 1

			// clamp to the last frame
			if si.idx == si.frames {
				si.idx = si.frames - 1
			}
		}

		// reset the timer
		si.timer = NewTimer(si.intervals[si.idx])
	}

	return next
}

func (si *SpriteIndex) PreviousFrame() *ebiten.Image {
	// get the current frame
	previous := si.images[si.idx]
	si.timer.Update()

	// if the timer is ready to advance to the previous frame
	if si.timer.IsReady() {
		if si.loop {
			// if looping, move to the previous frame
			si.idx = si.idx - 1
			if si.idx < 0 {
				si.idx += si.frames
			}
		} else {
			// clamp to the first frame
			si.idx = max(si.idx-1, 0)
		}

		// reset the timer
		si.timer = NewTimer(si.intervals[si.idx])
	}

	return previous
}

type DirectionalSpriteIndex struct {
	Up    *SpriteIndex
	Down  *SpriteIndex
	Left  *SpriteIndex
	Right *SpriteIndex
}

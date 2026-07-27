package engine

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteIndex struct {
	images         []*ebiten.Image
	frames         int
	intervals      []time.Duration
	idx            int
	timer          *Timer
	loop           bool
	duration       time.Duration
	width          int
	height         int
	skipFirstFrame bool
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

	width := rects[0].Dx()
	height := rects[0].Dy()

	return &SpriteIndex{
		images:    slice,
		frames:    len(rects),
		intervals: intervals,
		idx:       0,
		timer:     NewTimer(intervals[0]),
		loop:      loop,
		duration:  duration,
		width:     width,
		height:    height,
	}
}

func NewHorizontalSpriteIndex(spriteSheet *ebiten.Image, width int, height int, frames int, loop bool, intervals []time.Duration, skipFirstFrame bool) *SpriteIndex {
	return NewHorizontalSpriteIndexWithOffset(spriteSheet, width, height, 0, 0, frames, loop, intervals, skipFirstFrame)
}

func NewHorizontalSpriteIndexWithOffset(spriteSheet *ebiten.Image, width int, height int, xOffset int, yOffset int, frames int, loop bool, intervals []time.Duration, skipFirstFrame bool) *SpriteIndex {
	slice := make([]*ebiten.Image, frames)
	for i := range frames {
		rect := image.Rect(xOffset+width*i, yOffset, xOffset+width*(i+1), yOffset+height)
		slice[i] = spriteSheet.SubImage(rect).(*ebiten.Image)
	}

	duration := 0 * time.Second
	for i := range intervals {
		duration += intervals[i]
	}

	return &SpriteIndex{
		images:         slice,
		frames:         frames,
		intervals:      intervals,
		idx:            0,
		timer:          NewTimer(intervals[0]),
		loop:           loop,
		duration:       duration,
		width:          width,
		height:         height,
		skipFirstFrame: skipFirstFrame,
	}
}

func (si *SpriteIndex) Dx() int {
	return si.width
}

func (si *SpriteIndex) Dy() int {
	return si.height
}

func (si *SpriteIndex) Frames() []*ebiten.Image {
	return si.images
}

func (si *SpriteIndex) Index() int {
	return si.idx
}

func (si *SpriteIndex) Reset() {
	si.idx = 0
}

func (si *SpriteIndex) FirstFrame() *ebiten.Image {
	return si.images[0]
}

func (si *SpriteIndex) NextFrame() *ebiten.Image {
	// skip the first frame if configured
	if si.skipFirstFrame && si.idx == 0 {
		si.idx = 1
		si.timer = NewTimer(si.intervals[si.idx])
	}

	// get the current frame
	next := si.images[si.idx]
	si.timer.Update()

	// if the timer is ready to advance to the next frame
	if si.timer.IsReady() {
		if si.loop {
			// move to the next frame
			si.idx = (si.idx + 1) % si.frames

			// skip the first frame if configured
			if si.skipFirstFrame && si.idx == 0 {
				si.idx = 1
			}

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
	// TODO handle skipping the first frame
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

type LpcSpriteIndex struct {
	WalkUp    *SpriteIndex
	WalkLeft  *SpriteIndex
	WalkDown  *SpriteIndex
	WalkRight *SpriteIndex
	HurtDeath *SpriteIndex
}

func UniformFrameIntervals(duration time.Duration, count int) []time.Duration {
	intervals := make([]time.Duration, count)
	for i := range intervals {
		intervals[i] = duration
	}

	return intervals
}

type SpriteSheetDescriptor struct {
	Animations []SpriteSheetAnimationDescriptor `json:"animations"`
}

type SpriteSheetAnimationDescriptor struct {
	Height         int    `json:"height"`
	Width          int    `json:"width"`
	Row            int    `json:"row"`
	Frames         int    `json:"frames"`
	Type           string `json:"type"`
	Facing         string `json:"facing"`
	Loop           bool   `json:"loop"`
	SkipFirstFrame bool   `json:"skipFirstFrame"`
	SpecificFrames []int  `json:"specificFrames,omitempty"`
}

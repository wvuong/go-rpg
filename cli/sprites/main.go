package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wvuong/gogame/assets"
	"github.com/wvuong/gogame/engine"
)

const (
	screenWidth  = 1280
	screenHeight = 512
)

func main() {
	fs := os.DirFS("assets")
	files, err := os.ReadDir("assets/lpc")
	if err != nil {
		panic("Could not read directory")
	}

	types := make([]string, 0)
	for _, file := range files {
		if file.Type().IsRegular() && strings.HasSuffix(file.Name(), "-spritesheet.png") {
			types = append(types, file.Name())
		}
	}

	jsonBytes, err := assets.MustLoadJson(fs, "lpc/spritesheet.json")
	if err != nil {
		panic("Could not load json")
	}

	var descriptor engine.SpriteSheetDescriptor
	err = json.Unmarshal(jsonBytes, &descriptor)
	if err != nil {
		panic("Could not unmarshall json")
	}

	spritesheets := make([]*spritesheet, 0)
	for _, t := range types {
		img := assets.MustLoadImage(fs, "lpc/"+t)

		spriteIndexes := make([]*engine.SpriteIndex, 0)
		spriteNames := make([]string, 0)
		lpcSpriteIndex := &engine.LpcSpriteIndex{}

		for _, animation := range descriptor.Animations {
			si := engine.NewHorizontalSpriteIndexWithOffset(img, animation.Width, animation.Height,
				0, animation.Row*animation.Height,
				animation.Frames, animation.Loop,
				engine.UniformFrameIntervals(150*time.Millisecond, animation.Frames),
				animation.SkipFirstFrame)
			spriteIndexes = append(spriteIndexes, si)
			spriteNames = append(spriteNames, animation.Type+" "+animation.Facing)
			log.Println(t, animation.Type, animation.Facing)

			switch animation.Type {
			case "WALK":
				switch animation.Facing {
				case "UP":
					lpcSpriteIndex.WalkUp = si
				case "DOWN":
					lpcSpriteIndex.WalkDown = si
				case "LEFT":
					lpcSpriteIndex.WalkLeft = si
				case "RIGHT":
					lpcSpriteIndex.WalkRight = si
				}
			case "DEATH":
				lpcSpriteIndex.HurtDeath = si
			}
		}

		log.Println("Sprite sheet", t, "loaded", len(spriteIndexes), "animations")
		log.Println("LPC", lpcSpriteIndex)

		ss := &spritesheet{
			spriteIndexes: spriteIndexes,
			spriteNames:   spriteNames,
		}

		spritesheets = append(spritesheets, ss)
	}

	log.Println("Total sprite sheets", len(spritesheets))

	g := game{
		spriteSheets:     spritesheets,
		spriteSheetIndex: 0,
		animationIndex:   0,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

type spritesheet struct {
	name          string
	spriteIndexes []*engine.SpriteIndex
	spriteNames   []string
}

type game struct {
	spriteSheets     []*spritesheet
	spriteSheetIndex int
	animationIndex   int
}

func (g *game) Update() error {
	reset := false
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.spriteSheetIndex = engine.CyclePrevious(g.spriteSheetIndex, g.spriteSheets)
		reset = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.spriteSheetIndex = engine.CycleNext(g.spriteSheetIndex, g.spriteSheets)
		reset = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		ss := g.spriteSheets[g.spriteSheetIndex]
		g.animationIndex = engine.CyclePrevious(g.animationIndex, ss.spriteIndexes)
		reset = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		ss := g.spriteSheets[g.spriteSheetIndex]
		g.animationIndex = engine.CycleNext(g.animationIndex, ss.spriteIndexes)
		reset = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		reset = true
	}

	if reset {
		ss := g.spriteSheets[g.spriteSheetIndex]
		ss.spriteIndexes[g.animationIndex].Reset()
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	ss := g.spriteSheets[g.spriteSheetIndex]
	si := ss.spriteIndexes[g.animationIndex]

	op := &ebiten.DrawImageOptions{}
	// draw in the middle of the screen
	centeredX := screenWidth/2 - (si.Dx() / 2)
	centeredY := screenHeight/2 - (si.Dy() / 2)
	op.GeoM.Translate(float64(centeredX), float64(centeredY))
	screen.DrawImage(si.NextFrame(), op)

	// draw border
	vector.StrokeRect(screen, float32(centeredX), float32(centeredY),
		float32(si.Dx()), float32(si.Dy()), 2, color.White, false)

	// draw all frames on the bottom
	for idx, frame := range si.Frames() {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(idx*frame.Bounds().Dx()), float64(screenHeight-frame.Bounds().Dy()))
		screen.DrawImage(frame, op)
	}

	// draw debug information
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Sheet %d, Animation %d", g.spriteSheetIndex, g.animationIndex), 0, 0)
	ebitenutil.DebugPrintAt(screen, ss.spriteNames[g.animationIndex], 0, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Frame: %d", si.Index()), 0, 40)
}

func (g *game) Layout(width, height int) (int, int) {
	return width, height
}

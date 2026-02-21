package scene

import (
	"bytes"
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wvuong/gogame/engine"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

type TitleScene struct {
	director *Director
	ui       *ebitenui.UI
}

func NewTitleScene(config engine.GameConfig, state *engine.GameState, director *Director) *TitleScene {
	font := DefaultFont()

	center := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(colornames.Darkslategray),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(50, 50),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(
				widget.DirectionVertical,
			),
		)),
	)

	a := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(96, 24),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(
				widget.DirectionVertical,
			),
		)),
	)

	worldButton := widget.NewButton(
		widget.ButtonOpts.TextFace(&font),
		widget.ButtonOpts.TextColor(&widget.ButtonTextColor{
			Idle:    colornames.Gainsboro,
			Hover:   colornames.Gainsboro,
			Pressed: colornames.Gainsboro,
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         DefaultNineSlice(colornames.Darkslategray),
			Hover:        DefaultNineSlice(Mix(colornames.Darkslategray, colornames.Mediumseagreen, 0.4)),
			Disabled:     DefaultNineSlice(Mix(colornames.Darkslategray, colornames.Gainsboro, 0.8)),
			Pressed:      PressedNineSlice(Mix(colornames.Darkslategray, colornames.Black, 0.4)),
			PressedHover: PressedNineSlice(Mix(colornames.Darkslategray, colornames.Black, 0.4)),
		}),
		widget.ButtonOpts.TextLabel("World"),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(180, 48),
		),
		widget.ButtonOpts.ClickedHandler(
			func(args *widget.ButtonClickedEventArgs) {
				director.SwitchToTileMap()
			},
		),
	)

	a.AddChild(worldButton)

	b := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(96, 24),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(
				widget.DirectionVertical,
			),
		)),
	)

	battleButton := widget.NewButton(
		widget.ButtonOpts.TextFace(&font),
		widget.ButtonOpts.TextColor(&widget.ButtonTextColor{
			Idle:    colornames.Gainsboro,
			Hover:   colornames.Gainsboro,
			Pressed: colornames.Gainsboro,
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:         DefaultNineSlice(colornames.Darkslategray),
			Hover:        DefaultNineSlice(Mix(colornames.Darkslategray, colornames.Mediumseagreen, 0.4)),
			Disabled:     DefaultNineSlice(Mix(colornames.Darkslategray, colornames.Gainsboro, 0.8)),
			Pressed:      PressedNineSlice(Mix(colornames.Darkslategray, colornames.Black, 0.4)),
			PressedHover: PressedNineSlice(Mix(colornames.Darkslategray, colornames.Black, 0.4)),
		}),
		widget.ButtonOpts.TextLabel("Battle"),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(180, 48),
		),
	)

	b.AddChild(battleButton)

	center.AddChild(a)
	center.AddChild(b)

	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(colornames.Gainsboro),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	root.AddChild(center)

	return &TitleScene{
		director: director,
		ui:       &ebitenui.UI{Container: root},
	}
}

func (ts *TitleScene) Update() {
	ts.ui.Update()
}

func (ts *TitleScene) Draw(screen *ebiten.Image) {
	ts.ui.Draw(screen)
}

func DefaultFont() text.Face {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{
		Source: s,
		Size:   20,
	}
}

func DefaultNineSlice(base color.Color) *image.NineSlice {
	var size float32 = 64
	var tiles float32 = 16
	var radius float32 = 8

	tile := size / tiles
	//facet := Mix(base, colornames.Gainsboro, 0.2)

	img := ebiten.NewImage(int(size), int(size))

	drawOp := &vector.DrawPathOptions{}
	drawOp.AntiAlias = true
	drawOp.ColorScale.ScaleWithColor(base)

	fillOp := &vector.FillOptions{}
	fillOp.FillRule = vector.FillRuleEvenOdd

	path := RoundedRectPath(0, tile, size, size-tile, radius, radius, radius, radius)
	vector.FillPath(img, path, fillOp, drawOp)
	path = RoundedRectPath(0, tile, size, size-tile*2, radius, radius, radius, radius)
	vector.FillPath(img, path, fillOp, drawOp)
	path = RoundedRectPath(tile, tile*2, size-tile*2, size-tile*4, radius, radius, radius, radius)
	vector.FillPath(img, path, fillOp, drawOp)
	return image.NewNineSliceBorder(img, int(tile*4))
}

func PressedNineSlice(base color.Color) *image.NineSlice {
	var size float32 = 64
	var tiles float32 = 16
	var radius float32 = 8

	tile := size / tiles
	//facet := Mix(base, colornames.Gainsboro, 0.2)

	img := ebiten.NewImage(int(size), int(size))

	drawOp := &vector.DrawPathOptions{}
	drawOp.AntiAlias = true
	drawOp.ColorScale.ScaleWithColor(base)

	fillOp := &vector.FillOptions{}
	fillOp.FillRule = vector.FillRuleEvenOdd

	path := RoundedRectPath(0, 0, size, size, radius, radius, radius, radius)
	vector.FillPath(img, path, fillOp, drawOp)
	path = RoundedRectPath(tile, tile, size-tile*2, size-tile*2, radius, radius, radius, radius)
	vector.FillPath(img, path, fillOp, drawOp)

	return image.NewNineSliceBorder(img, int(tile*4))
}

func Mix(a, b color.Color, percent float64) color.Color {
	rgba := func(c color.Color) (r, g, b, a uint8) {
		r16, g16, b16, a16 := c.RGBA()
		return uint8(r16 >> 8), uint8(g16 >> 8), uint8(b16 >> 8), uint8(a16 >> 8)
	}
	lerp := func(x, y uint8) uint8 {
		return uint8(math.Round(float64(x) + percent*(float64(y)-float64(x))))
	}
	r1, g1, b1, a1 := rgba(a)
	r2, g2, b2, a2 := rgba(b)

	return color.RGBA{
		R: lerp(r1, r2),
		G: lerp(g1, g2),
		B: lerp(b1, b2),
		A: lerp(a1, a2),
	}
}

func RoundedRectPath(x, y, w, h, tl, tr, br, bl float32) *vector.Path {
	path := &vector.Path{}

	path.Arc(x+w-tr, y+tr, tr, 3*math.Pi/2, 0, vector.Clockwise)
	path.LineTo(x+w, y+h-br)
	path.Arc(x+w-br, y+h-br, br, 0, math.Pi/2, vector.Clockwise)
	path.LineTo(x+bl, y+h)
	path.Arc(x+bl, y+h-bl, bl, math.Pi/2, math.Pi, vector.Clockwise)
	path.LineTo(x, y+tl)
	path.Arc(x+tl, y+tl, tl, math.Pi, 3*math.Pi/2, vector.Clockwise)
	path.Close()

	return path
}

package engine

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

type TileMap struct {
	layers     [][]int
	cols       int
	rows       int
	tileSize   int
	mapWidth   int
	mapHeight  int
	tileSheets []*ebiten.Image
}

func NewTileMap(tileSheets []*ebiten.Image, layers [][]int, cols int, rows int, tileSize int) *TileMap {
	return &TileMap{
		layers:     layers,
		cols:       cols,
		rows:       rows,
		tileSize:   tileSize,
		mapWidth:   cols * tileSize,
		mapHeight:  rows * tileSize,
		tileSheets: tileSheets,
	}
}

func (m *TileMap) getTile(layer int, col int, row int) int {
	if col >= m.cols || row >= m.rows {
		return 0
	}

	return m.layers[layer][row*m.cols+col]
}

func (m *TileMap) isSolidTileAtXY(x float64, y float64) bool {
	// right 680 + 24 = 704
	// bottom 96 + 32 = 128
	col := m.getCol(x)
	row := m.getRow(y)

	tile := m.getTile(0, col, row)
	return tile > 2
}

func (m *TileMap) getCol(x float64) int {
	// right 704/64 = 11.0
	// bottom 128/64 = 2.0
	return int(math.Floor(x / float64(m.tileSize)))
}

func (m *TileMap) getRow(y float64) int {
	return int(math.Floor(y / float64(m.tileSize)))
}

func (m *TileMap) getX(col int) float64 {
	return float64(col * m.tileSize)
}

func (m *TileMap) getY(row int) float64 {
	return float64(row * m.tileSize)
}

func (m *TileMap) drawTile(screen *ebiten.Image, layer int, col int, row int, x float64, y float64, debug *Debug) {
	// this is the raw tile index from the map data
	tileId := m.getTile(layer, col, row)

	// 0 => empty tile
	if tileId != 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)

		tileIdx := tileId - 1 // subtract 1 to get the tile array index

		// create a rect of the tile in the tilesheet
		sx := tileIdx * m.tileSize
		rect := image.Rect(sx, 0, sx+m.tileSize, m.tileSize)

		// if the tile is solid, tint it red
		if debug.Enabled && tileId > 2 {
			op.ColorScale.ScaleWithColor(colornames.Red)
		}

		// draw the tile to the screen
		// ebiten will automatically clip the tile if it's partially off-screen
		screen.DrawImage(m.tileSheets[layer].SubImage(rect).(*ebiten.Image), op)

		if debug.Enabled {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("(%d)", tileId),
				int(x+float64(m.tileSize)/2), int(y+float64(m.tileSize)/2))
		}
	}
}

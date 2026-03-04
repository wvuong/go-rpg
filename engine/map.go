package engine

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TileMap struct {
	layers     [][]int
	impassable []int
	cols       int
	rows       int
	tileSize   int
	mapWidth   int
	mapHeight  int
	tileImages []*ebiten.Image
}

func NewTileMap(tileImages []*ebiten.Image, layers [][]int, impassable []int, cols int, rows int, tileSize int) *TileMap {
	return &TileMap{
		layers:     layers,
		impassable: impassable,
		cols:       cols,
		rows:       rows,
		tileSize:   tileSize,
		mapWidth:   cols * tileSize,
		mapHeight:  rows * tileSize,
		tileImages: tileImages,
	}
}

func (m *TileMap) getTile(layer int, col int, row int) int {
	if col >= m.cols || row >= m.rows {
		return 0
	}

	return m.layers[layer][row*m.cols+col]
}

// called by player
func (m *TileMap) isSolidTileAtXY(x float64, y float64) bool {
	col := m.getCol(x)
	row := m.getRow(y)

	idx := row*m.cols + col
	return m.impassable[idx] == 1
}

// called by player
func (m *TileMap) getCol(x float64) int {
	return int(math.Floor(x / float64(m.tileSize)))
}

// called by player
func (m *TileMap) getRow(y float64) int {
	return int(math.Floor(y / float64(m.tileSize)))
}

// called by player
func (m *TileMap) getX(col int) float64 {
	return float64(col * m.tileSize)
}

// called by player
func (m *TileMap) getY(row int) float64 {
	return float64(row * m.tileSize)
}

// called by camera
func (m *TileMap) drawTile(screen *ebiten.Image, layer int, col int, row int, x float64, y float64, debug *Debug) {
	// this is the raw tile index from the map data
	tileId := m.getTile(layer, col, row)

	// 0 => empty tile
	if tileId != 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)

		// draw the tile to the screen
		// ebiten will automatically clip the tile if it's partially off-screen
		screen.DrawImage(m.tileImages[tileId], op)

		if debug.Enabled {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("(%d)", tileId),
				int(x+float64(m.tileSize)/2), int(y+float64(m.tileSize)/2))
		}
	}
}

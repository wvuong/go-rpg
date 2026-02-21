package engine

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type TileMap struct {
	Layers    [][]int
	Cols      int
	Rows      int
	TileSize  int
	MapWidth  int
	MapHeight int
	TileSheet *ebiten.Image
}

func NewMap(tileSheet *ebiten.Image, layers [][]int, cols int, rows int, tileSize int) *TileMap {
	return &TileMap{
		Layers:    layers,
		Cols:      cols,
		Rows:      rows,
		TileSize:  tileSize,
		MapWidth:  cols * tileSize,
		MapHeight: rows * tileSize,
		TileSheet: tileSheet,
	}
}

func (m *TileMap) GetTile(layer int, col int, row int) int {
	if col >= m.Cols || row >= m.Rows {
		return 0
	}

	return m.Layers[layer][row*m.Cols+col]
}

func (m *TileMap) isSolidTileAtXY(x float64, y float64) bool {
	// right 680 + 24 = 704
	// bottom 96 + 32 = 128
	col := m.getCol(x)
	row := m.getRow(y)

	tile := m.GetTile(0, col, row)
	return tile > 2
}

func (m *TileMap) getCol(x float64) int {
	// right 704/64 = 11.0
	// bottom 128/64 = 2.0
	return int(math.Floor(x / float64(m.TileSize)))
}

func (m *TileMap) getRow(y float64) int {
	return int(math.Floor(y / float64(m.TileSize)))
}

func (m *TileMap) getX(col int) float64 {
	return float64(col * m.TileSize)
}

func (m *TileMap) getY(row int) float64 {
	return float64(row * m.TileSize)
}

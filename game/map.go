package game

import "math"

type TileMap struct {
	Layers   [][]int
	Cols     int
	Rows     int
	TileSize int
}

func NewMap(layers [][]int, cols int, rows int, tileSize int) *TileMap {
	return &TileMap{Layers: layers, Cols: cols, Rows: rows, TileSize: tileSize}
}

func (m *TileMap) GetTile(layer int, col int, row int) int {
	if col >= m.Cols || row >= m.Rows {
		return 0
	}

	return m.Layers[layer][row*m.Cols+col]
}

func (m *TileMap) isSolidTileAtXY(col int, row int) bool {
	tile := m.GetTile(0, col, row)
	return tile > 2
}

func (m *TileMap) getCol(x int) int {
	return int(math.Floor(float64(x) / float64(m.TileSize)))
}

func (m *TileMap) getRow(y int) int {
	return int(math.Floor(float64(y) / float64(m.TileSize)))
}

func (m *TileMap) getX(col int) float64 {
	return float64(col * m.TileSize)
}

func (m *TileMap) getY(row int) float64 {
	return float64(row * m.TileSize)
}

package game

type TileMap struct {
	Layers [][]int
	Cols   int
	Rows   int
}

func NewMap(layers [][]int, cols int, rows int) *TileMap {
	return &TileMap{Layers: layers, Cols: cols, Rows: rows}
}

func (m *TileMap) GetTile(layer int, col int, row int) int {
	if col >= m.Cols || row >= m.Rows {
		return 0
	}

	return m.Layers[layer][row*m.Cols+col]
}

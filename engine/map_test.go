package engine

import "testing"

// newTestTileMap mirrors the geometry of assets/levels/world.tmx: 30x20 tiles at 64px.
func newTestTileMap() *TileMap {
	cols, rows, tileSize := 30, 20, 64
	return NewTileMap(nil, [][]int{make([]int, cols*rows)}, make([]int, cols*rows), cols, rows, tileSize)
}

// Out-of-bounds coordinates used to index the flattened collision mask directly,
// which either panicked or wrapped onto a neighbouring row.
func TestIsSolidTileAtXYOutsideMapIsSolid(t *testing.T) {
	m := newTestTileMap()

	cases := []struct {
		name string
		x, y float64
	}{
		{"left of map", -24, 30},    // panicked: index out of range [-1]
		{"above map", 100, -1},      // panicked: negative index
		{"below map", 100, 1311},    // panicked: index out of range [601]
		{"right of map", 1943, 100}, // silently wrapped to the next row
		{"far outside", 999999, 999999},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !m.isSolidTileAtXY(c.x, c.y) {
				t.Errorf("isSolidTileAtXY(%v, %v) = false, want true so the player is stopped at the map edge", c.x, c.y)
			}
		})
	}
}

func TestIsSolidTileAtXYInsideMapUsesMask(t *testing.T) {
	m := newTestTileMap()

	// mark the tile at col 5, row 3 impassable
	m.impassable[3*m.cols+5] = 1

	if !m.isSolidTileAtXY(5*64, 3*64) {
		t.Error("tile marked impassable reported as passable")
	}
	if m.isSolidTileAtXY(6*64, 3*64) {
		t.Error("tile not marked impassable reported as solid")
	}
}

// getTile guarded the upper bounds but not negatives.
func TestGetTileOutsideMapIsEmpty(t *testing.T) {
	m := newTestTileMap()

	for _, c := range []struct{ col, row int }{
		{-1, 0}, {0, -1}, {-1, -1}, {30, 0}, {0, 20}, {30, 20},
	} {
		if got := m.getTile(0, c.col, c.row); got != 0 {
			t.Errorf("getTile(0, %d, %d) = %d, want 0 (empty)", c.col, c.row, got)
		}
	}
}

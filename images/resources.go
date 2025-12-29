package images

import (
	_ "embed"
)

var (
	//go:embed tiles.png
	Tiles_png []byte

	//go:embed universal-lpc-sprite_male_01_walk-3frame.png
	Sprite_png []byte
)

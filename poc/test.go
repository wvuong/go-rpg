package main

import (
	"image"
	"log"
)

const (
	cols     = 8
	rows     = 8
	tileSize = 64
)

func main() {
	layer := []int{
		1, 3, 3, 3, 1, 1, 3, 1, // 0
		1, 1, 1, 1, 1, 1, 1, 1, // 1
		1, 1, 1, 1, 1, 2, 1, 1, // 2
		1, 1, 1, 1, 1, 1, 1, 1, // 3
		1, 1, 1, 2, 1, 1, 1, 1, // 4
		1, 1, 1, 1, 2, 1, 1, 1, // 5
		1, 1, 1, 1, 2, 1, 1, 1, // 6
		1, 1, 1, 0, 0, 1, 1, 1, // 7
	}
	//  0. 1. 2. 3. 4. 5. 6. 7.

	for c := range cols {
		for r := range rows {
			tile := GetTile(layer, c, r)
			if tile != 0 {
				// 0 => empty tile
				tx := float64(c * tileSize)
				ty := float64(r * tileSize)
				log.Println("c", c, "r:", r, "tile:", tile, "tx:", tx, "ty:", ty)
			}
		}
	}

	for i := range 5 {
		sx := i * tileSize
		rect := image.Rect(sx, 0, sx+tileSize, tileSize)
		log.Println(i, rect)
	}
}

func GetTile(layer []int, col int, row int) int {
	return layer[row*cols+col]
}

package assets

import (
	"embed"
	_ "embed"
	"fmt"
	"image"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/wvuong/gogame/engine"
)

var (
	//go:embed *
	assetsFS embed.FS

	Sprite_png *ebiten.Image

	TileMaps map[string]*engine.TileMap
)

func MustLoadAssets() {
	Sprite_png = mustLoadImage("universal-lpc-sprite_male_01_walk-3frame.png")

	levelLoader := newTileMapLoader()
	TileMaps = levelLoader.MustLoadTileMaps()
}

func mustLoadImage(name string) *ebiten.Image {
	f, err := assetsFS.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

type tileMapLoader struct {
}

func newTileMapLoader() *tileMapLoader {
	return &tileMapLoader{}
}

func (l *tileMapLoader) MustLoadTileMaps() map[string]*engine.TileMap {
	paths, err := fs.Glob(assetsFS, filepath.Join("levels", "world*.tmx"))
	if err != nil {
		panic(err)
	}

	tileMaps := make(map[string]*engine.TileMap)
	for _, path := range paths {
		tileMap := l.MustLoadTileMap(path)
		tileMaps[path] = tileMap
		log.Println("Loaded tilemap", path)
	}

	return tileMaps
}

func (l *tileMapLoader) MustLoadTileMap(path string) *engine.TileMap {
	tileMap, err := tiled.LoadFile(path, tiled.WithFileSystem(assetsFS))
	if err != nil {
		panic(err)
	}

	log.Println(tileMap.Width, tileMap.Height)
	log.Println(tileMap.TileWidth, tileMap.TileHeight)
	tilesets := tileMap.Tilesets
	log.Println("tilesets", len(tilesets))

	tiles := 0
	for _, tileset := range tilesets {
		tiles += tileset.TileCount
	}

	log.Println("total tiles", tiles)

	// Store sub-image references for each tile contained in the tileset. A
	// sub-image reference only isolates part of a source image. Thus, each
	// of the sub-image references point to the tileset image we just loaded.
	tileImages := make([]*ebiten.Image, tiles+1) // +1 because tile IDs are 1-based, so we need an extra slot for the 0 ID (empty tile).

	for tilesetIndex, tileset := range tilesets {
		log.Println("tileset", "#"+strconv.Itoa(tilesetIndex), tileset.FirstGID, tileset.Name, tileset.Image.Source, tileset.TileCount, len(tileset.Tiles))
		img := mustLoadImage(filepath.Join("levels", tileset.Image.Source))

		for i := uint32(0); i < uint32(tileset.TileCount); i++ {
			// Calculate sub-image rect for this tile.
			globalTileId := tileset.FirstGID + i
			r := tileset.GetTileRect(i)

			// Store sub-image reference for this tile.
			tileImages[globalTileId] = img.SubImage(r).(*ebiten.Image)
			log.Println("["+strconv.Itoa(int(globalTileId))+"]", "<--", tileset.FirstGID, i, r)
		}

		for tileIndex, tile := range tileset.Tiles {
			log.Println("tileset tile", "#"+strconv.Itoa(tileIndex), tile.ID, tile.Properties.GetBool("terrain-impassable"))
		}
	}

	layers := make([][]int, len(tileMap.Layers))
	impassable := make([]int, tileMap.Width*tileMap.Height)

	for layerIndex, layer := range tileMap.Layers {
		layerData := make([]int, tileMap.Width*tileMap.Height)

		log.Println("layer", "#"+strconv.Itoa(layerIndex), layer.ID, layer.Name, len(layer.Tiles))
		for idy := range tileMap.Height {
			for idx := range tileMap.Width {
				arrayIndex := idx + (idy * tileMap.Width)
				layerTile := layer.Tiles[arrayIndex]

				if layerTile.Tileset != nil {
					gid := layerTile.ID + layerTile.Tileset.FirstGID
					layerData[arrayIndex] = int(gid)

					if layerIndex != 0 {
						// Check if the tile has the "terrain-impassable" property set to true.
						tileSetTile, err := layerTile.Tileset.GetTilesetTile(layerTile.ID)
						if err == nil {
							val := tileSetTile.Properties.GetBool("terrain-impassable")
							if val {
								impassable[arrayIndex] = 1
							}
						}
					}

					fmt.Printf("%02d ", gid)

				} else {
					layerData[arrayIndex] = 0
					fmt.Printf("-- ")
				}
			}
			fmt.Println()
		}

		layers[layerIndex] = layerData
	}

	log.Println("impassable height map")
	for idy := range tileMap.Height {
		for idx := range tileMap.Width {
			arrayIndex := idx + (idy * tileMap.Width)
			fmt.Printf("%d ", impassable[arrayIndex])
		}
		fmt.Println()
	}

	return engine.NewTileMap(tileImages, layers, impassable, tileMap.Width, tileMap.Height, tileMap.TileWidth)
}

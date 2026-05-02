package assets

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/wvuong/gogame/engine"
)

var (
	//go:embed *
	assetsFS embed.FS

	Sprite_png *ebiten.Image

	BlackWarrior_Idle   *ebiten.Image
	BlackWarrior_Run    *ebiten.Image
	BlackWarrior_Attack *ebiten.Image
	BlackWarrior_Guard  *ebiten.Image

	BlueWarrior_Idle   *ebiten.Image
	BlueWarrior_Run    *ebiten.Image
	BlueWarrior_Attack *ebiten.Image
	BlueWarrior_Guard  *ebiten.Image

	BlackArcher_Idle   *ebiten.Image
	BlackArcher_Run    *ebiten.Image
	BlackArcher_Attack *ebiten.Image
	Arrow              *ebiten.Image

	BlackLancer_Idle   *ebiten.Image
	BlackLancer_Run    *ebiten.Image
	BlackLancer_Attack *ebiten.Image
	BlackLancer_Guard  *ebiten.Image

	LpcSpriteNames []string
	LpcSpritesMap  map[string]*engine.LpcSpriteIndex

	TileMaps map[string]*engine.TileMap
)

func MustLoadAssets() {
	Sprite_png = mustLoadImage("universal-lpc-sprite_male_01_walk-3frame.png")

	BlackWarrior_Idle = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Warrior/Warrior_Idle.png")
	BlackWarrior_Run = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Warrior/Warrior_Run.png")
	BlackWarrior_Attack = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Warrior/Warrior_Attack1.png")
	BlackWarrior_Guard = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Warrior/Warrior_Guard.png")

	BlueWarrior_Idle = mustLoadImage("Tiny Swords (Free Pack)/Units/Blue Units/Warrior/Warrior_Idle.png")
	BlueWarrior_Run = mustLoadImage("Tiny Swords (Free Pack)/Units/Blue Units/Warrior/Warrior_Run.png")
	BlueWarrior_Attack = mustLoadImage("Tiny Swords (Free Pack)/Units/Blue Units/Warrior/Warrior_Attack1.png")
	BlueWarrior_Guard = mustLoadImage("Tiny Swords (Free Pack)/Units/Blue Units/Warrior/Warrior_Guard.png")

	BlackArcher_Idle = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Archer/Archer_Idle.png")
	BlackArcher_Run = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Archer/Archer_Run.png")
	BlackArcher_Attack = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Archer/Archer_Shoot.png")
	Arrow = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Archer/Arrow.png")

	BlackLancer_Idle = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Lancer/Lancer_Idle.png")
	BlackLancer_Run = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Lancer/Lancer_Run.png")
	BlackLancer_Attack = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Lancer/Lancer_Right_Attack.png")
	BlackLancer_Guard = mustLoadImage("Tiny Swords (Free Pack)/Units/Black Units/Lancer/Lancer_Right_Defence.png")

	levelLoader := newTileMapLoader()
	TileMaps = levelLoader.MustLoadTileMaps()

	mustLoadLpcSprites()
}

func mustLoadImage(name string) *ebiten.Image {
	return MustLoadImage(assetsFS, name)
}

func MustLoadImage(filesystem fs.FS, name string) *ebiten.Image {
	f, err := filesystem.Open(name)
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

func MustLoadJson(filesystem fs.FS, name string) ([]byte, error) {
	return fs.ReadFile(filesystem, name)
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

func mustLoadLpcSprites() {
	fs := os.DirFS("assets")
	files, err := os.ReadDir("assets/lpc")
	if err != nil {
		panic("Could not read directory")
	}

	fileNames := make([]string, 0)
	LpcSpriteNames = make([]string, 0)
	for _, file := range files {
		if file.Type().IsRegular() && strings.HasSuffix(file.Name(), "-spritesheet.png") {
			fileNames = append(fileNames, file.Name())
			name := strings.TrimSuffix(file.Name(), "-spritesheet.png")
			LpcSpriteNames = append(LpcSpriteNames, name)
			log.Println("Found LPC sprite sheet", file.Name())
		}
	}

	jsonBytes, err := MustLoadJson(fs, "lpc/spritesheet.json")
	if err != nil {
		panic("Could not load json")
	}

	var descriptor engine.SpriteSheetDescriptor
	err = json.Unmarshal(jsonBytes, &descriptor)
	if err != nil {
		panic("Could not unmarshall json")
	}

	LpcSpritesMap := make(map[string]*engine.LpcSpriteIndex)

	for _, t := range fileNames {
		img := MustLoadImage(fs, "lpc/"+t)

		spriteIndexes := make([]*engine.SpriteIndex, 0)
		spriteNames := make([]string, 0)
		lpcSpriteIndex := &engine.LpcSpriteIndex{}

		for _, animation := range descriptor.Animations {
			si := engine.NewHorizontalSpriteIndexWithOffset(img, animation.Width, animation.Height,
				0, animation.Row*animation.Height,
				animation.Frames, animation.Loop,
				engine.UniformFrameIntervals(100*time.Millisecond, animation.Frames),
				animation.SkipFirstFrame)
			spriteIndexes = append(spriteIndexes, si)
			spriteNames = append(spriteNames, animation.Type+" "+animation.Facing)
			log.Println(t, animation.Type, animation.Facing)

			switch animation.Type {
			case "WALK":
				switch animation.Facing {
				case "UP":
					lpcSpriteIndex.WalkUp = si
				case "DOWN":
					lpcSpriteIndex.WalkDown = si
				case "LEFT":
					lpcSpriteIndex.WalkLeft = si
				case "RIGHT":
					lpcSpriteIndex.WalkRight = si
				}
			case "DEATH":
				lpcSpriteIndex.HurtDeath = si
			}
		}

		LpcSpritesMap[t] = lpcSpriteIndex
		log.Println("Sprite sheet", t, "loaded", len(spriteIndexes), "animations")
	}
}

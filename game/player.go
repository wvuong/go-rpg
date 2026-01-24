package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	TileMap *TileMap
	Sprite  *Sprite

	DirectionalSpriteIndex *DirectionalSpriteIndex

	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

func NewPlayer(tileMap *TileMap, directionalSpriteIndex *DirectionalSpriteIndex, x float64, y float64) *Player {
	startingImage := directionalSpriteIndex.Down.NextFrame()
	sprite := NewSprite(startingImage, x, y)

	return &Player{
		TileMap:                tileMap,
		Sprite:                 sprite,
		DirectionalSpriteIndex: directionalSpriteIndex,
	}
}

func (p *Player) Update() {
	// move player with arrow keys input
	var dirX, dirY float64
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dirY = -playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dirY = playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dirX = -playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dirX = playerSpeed
	}

	// update player position
	p.Sprite.Position.X += dirX
	p.Sprite.Position.Y += dirY

	// check for collisions at new player position
	// calculate the player's bounding box
	left := p.Sprite.Position.X - float64(p.Sprite.Dx)/2
	right := p.Sprite.Position.X + float64(p.Sprite.Dx)/2 - 1
	top := p.Sprite.Position.Y - float64(p.Sprite.Dy)/2
	bottom := p.Sprite.Position.Y + float64(p.Sprite.Dy)/2 - 1

	// store the bounding box values
	p.Left = left
	p.Right = right
	p.Top = top
	p.Bottom = bottom

	// check for collisions with solid tiles on the bounding box
	collision := p.TileMap.isSolidTileAtXY(left, top) ||
		p.TileMap.isSolidTileAtXY(right, top) ||
		p.TileMap.isSolidTileAtXY(right, bottom) ||
		p.TileMap.isSolidTileAtXY(left, bottom)

	// if collision, reset player position based on movement direction
	if collision {
		if dirY > 0 {
			// moving down
			row := p.TileMap.getRow(bottom)
			// align player to top edge of tile
			p.Sprite.Position.Y = -float64(p.Sprite.Dy)/2 + p.TileMap.getY(row)
		} else if dirY < 0 {
			// moving up
			row := p.TileMap.getRow(top)
			// align player to bottom edge of tile
			p.Sprite.Position.Y = float64(p.Sprite.Dy)/2 + p.TileMap.getY(row+1)
		} else if dirX > 0 {
			// moving right
			col := p.TileMap.getCol(right)
			// align player to left edge of tile
			p.Sprite.Position.X = -float64(p.Sprite.Dx)/2 + p.TileMap.getX(col)
		} else if dirX < 0 {
			// moving left
			col := p.TileMap.getCol(left)
			// align player to right edge of tile
			p.Sprite.Position.X = float64(p.Sprite.Dx)/2 + p.TileMap.getX(col+1)
		}
	}

	// clamp player position to map bounds
	x := math.Max(0, math.Min(p.Sprite.Position.X, mapWidth))
	y := math.Max(0, math.Min(p.Sprite.Position.Y, mapHeight))
	p.Sprite.Position.X = x
	p.Sprite.Position.Y = y

	// update animation based on movement direction
	if dirY < 0 {
		// moving up
		p.Sprite.Image = p.DirectionalSpriteIndex.Up.NextFrame()

	} else if dirY > 0 {
		// moving down
		p.Sprite.Image = p.DirectionalSpriteIndex.Down.NextFrame()

	} else if dirX < 0 {
		// moving left
		p.Sprite.Image = p.DirectionalSpriteIndex.Left.NextFrame()

	} else if dirX > 0 {
		// moving right
		p.Sprite.Image = p.DirectionalSpriteIndex.Right.NextFrame()
	}
}

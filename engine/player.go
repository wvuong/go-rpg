package engine

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	tileMap  *TileMap
	Sprite   *Sprite
	Position *Vector

	directionalSpriteIndex *DirectionalSpriteIndex
	speed                  float64

	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

func NewPlayer(tileMap *TileMap, directionalSpriteIndex *DirectionalSpriteIndex, position *Vector, speed float64) *Player {
	startingImage := directionalSpriteIndex.Down.NextFrame()
	sprite := NewSprite(startingImage)

	return &Player{
		tileMap:                tileMap,
		Sprite:                 sprite,
		Position:               position,
		directionalSpriteIndex: directionalSpriteIndex,
		speed:                  speed,
	}
}

func (p *Player) Update() {
	// move player with arrow keys input
	var dirX, dirY float64
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dirY = -p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dirY = p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dirX = -p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dirX = p.speed
	}

	// update player position
	p.Position.X += dirX
	p.Position.Y += dirY

	// check for collisions at new player position
	// calculate the player's bounding box
	left := p.Position.X - float64(p.Sprite.Dx)/2
	right := p.Position.X + float64(p.Sprite.Dx)/2 - 1
	top := p.Position.Y - float64(p.Sprite.Dy)/2
	bottom := p.Position.Y + float64(p.Sprite.Dy)/2 - 1

	// store the bounding box values
	p.Left = left
	p.Right = right
	p.Top = top
	p.Bottom = bottom

	// check for collisions with solid tiles on the bounding box
	collision := p.tileMap.isSolidTileAtXY(left, top) ||
		p.tileMap.isSolidTileAtXY(right, top) ||
		p.tileMap.isSolidTileAtXY(right, bottom) ||
		p.tileMap.isSolidTileAtXY(left, bottom)

	// if collision, reset player position based on movement direction
	if collision {
		if dirY > 0 {
			// moving down
			row := p.tileMap.getRow(bottom)
			// align player to top edge of tile
			p.Position.Y = -float64(p.Sprite.Dy)/2 + p.tileMap.getY(row)
		} else if dirY < 0 {
			// moving up
			row := p.tileMap.getRow(top)
			// align player to bottom edge of tile
			p.Position.Y = float64(p.Sprite.Dy)/2 + p.tileMap.getY(row+1)
		} else if dirX > 0 {
			// moving right
			col := p.tileMap.getCol(right)
			// align player to left edge of tile
			p.Position.X = -float64(p.Sprite.Dx)/2 + p.tileMap.getX(col)
		} else if dirX < 0 {
			// moving left
			col := p.tileMap.getCol(left)
			// align player to right edge of tile
			p.Position.X = float64(p.Sprite.Dx)/2 + p.tileMap.getX(col+1)
		}
	}

	// clamp player position to map bounds
	x := math.Max(0, math.Min(p.Position.X, float64(p.tileMap.mapWidth)))
	y := math.Max(0, math.Min(p.Position.Y, float64(p.tileMap.mapHeight)))
	p.Position.X = x
	p.Position.Y = y

	// update animation based on movement direction
	if dirY < 0 {
		// moving up
		p.Sprite.Image = p.directionalSpriteIndex.Up.NextFrame()

	} else if dirY > 0 {
		// moving down
		p.Sprite.Image = p.directionalSpriteIndex.Down.NextFrame()
	} else if dirX < 0 {
		// moving left
		p.Sprite.Image = p.directionalSpriteIndex.Left.NextFrame()

	} else if dirX > 0 {
		// moving right
		p.Sprite.Image = p.directionalSpriteIndex.Right.NextFrame()
	}
}

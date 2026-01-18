package game

type Player struct {
	Sprite *Sprite
	Left   int
	Right  int
	Top    int
	Bottom int
}

func NewPlayer(sprite *Sprite) *Player {
	return &Player{
		Sprite: sprite,
	}
}

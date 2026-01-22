package game

type Player struct {
	Sprite *Sprite
	Left   float64
	Right  float64
	Top    float64
	Bottom float64
}

func NewPlayer(sprite *Sprite) *Player {
	return &Player{
		Sprite: sprite,
	}
}

package engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type BattleAnimations struct {
	Idle      *SpriteIndex
	Run       *SpriteIndex
	AttackOne *SpriteIndex
	AttackTwo *SpriteIndex
	Heal      *SpriteIndex
	Guard     *SpriteIndex
}

type BattlerState int

const (
	BattlerIdle BattlerState = iota
	BattlerRunning
	BattlerAttacking
	BattlerGuarding
)

type Battler struct {
	sprite             *Sprite
	battleAnimations   *BattleAnimations
	forwardAnimation   bool
	attacked           bool
	state              BattlerState
	bottomMessage      string
	bottomMessageTimer *Timer
	spriteIndex        *SpriteIndex
}

func NewBattler(sprite *Sprite, battleAnimations *BattleAnimations, position *Vector) *Battler {
	// set the initial position of the sprite
	sprite.ScreenPosition = position

	return &Battler{
		sprite:           sprite,
		battleAnimations: battleAnimations,
		forwardAnimation: true,
		state:            BattlerIdle,
		spriteIndex:      battleAnimations.Idle,
	}
}

func (b *Battler) Position() *Vector {
	return b.sprite.ScreenPosition
}

func (b *Battler) Dx() int {
	return b.sprite.Dx
}

func (b *Battler) UseForwardAnimation(forwardAnimation bool) {
	b.forwardAnimation = forwardAnimation
}

func (b *Battler) MoveY(delta float64) {
	b.sprite.ScreenPosition.Y += delta
}

func (b *Battler) MoveX(delta float64) {
	b.sprite.ScreenPosition.X += delta
}

func (b *Battler) SetState(state BattlerState) {
	b.spriteIndex.Reset()
	b.state = state
}

func (b *Battler) SetBottomMessage(message string, duration time.Duration) {
	b.bottomMessage = message
	b.bottomMessageTimer = NewTimer(duration)
}

func (b *Battler) Attacked(attacked bool) {
	b.attacked = attacked
}

func (b *Battler) Update() {
	if b.bottomMessageTimer != nil {
		b.bottomMessageTimer.Update()
		if b.bottomMessageTimer.IsReady() {
			b.bottomMessage = ""
			b.bottomMessageTimer = nil
		}
	}

	switch b.state {
	case BattlerIdle:
		if b.battleAnimations.Idle != nil {
			b.spriteIndex = b.battleAnimations.Idle
			if b.forwardAnimation {
				b.sprite.Image = b.battleAnimations.Idle.NextFrame()
			} else {
				b.sprite.Image = b.battleAnimations.Idle.PreviousFrame()
			}
		}
	case BattlerRunning:
		if b.battleAnimations.Run != nil {
			b.spriteIndex = b.battleAnimations.Run
			if b.forwardAnimation {
				b.sprite.Image = b.battleAnimations.Run.NextFrame()
			} else {
				b.sprite.Image = b.battleAnimations.Run.PreviousFrame()
			}
		}
	case BattlerAttacking:
		if b.battleAnimations.AttackOne != nil {
			b.spriteIndex = b.battleAnimations.AttackOne
			if b.forwardAnimation {
				b.sprite.Image = b.battleAnimations.AttackOne.NextFrame()
			} else {
				b.sprite.Image = b.battleAnimations.AttackOne.PreviousFrame()
			}
		}
	case BattlerGuarding:
		if b.battleAnimations.Guard != nil {
			b.spriteIndex = b.battleAnimations.Guard
			if b.forwardAnimation {
				b.sprite.Image = b.battleAnimations.Guard.NextFrame()
			} else {
				b.sprite.Image = b.battleAnimations.Guard.PreviousFrame()
			}
		}
	}
}

func (b *Battler) Draw(screen *ebiten.Image, debug *Debug) {
	op := &ebiten.DrawImageOptions{}
	if b.attacked {
		// if the battler is currently being attacked, we want to tint it red
		cs := ebiten.ColorScale{}
		cs.Scale(1.5, 0.5, 0.5, 1)
		op.ColorScale = cs
	}

	b.sprite.Draw(screen, op, debug)

	// if the battler has a bottom message, draw it below the battler
	if b.bottomMessage != "" {
		halfW, halfH := float64(b.sprite.Dx)/2, float64(b.sprite.Dy)/2
		centeredX := b.sprite.ScreenPosition.X - halfW
		centeredY := b.sprite.ScreenPosition.Y - halfH
		ebitenutil.DebugPrintAt(screen, b.bottomMessage, int(centeredX)+(b.sprite.Dx/2), int(centeredY)+b.sprite.Dy)
	}
}

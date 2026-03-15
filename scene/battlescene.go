package scene

import (
	"math"
	"slices"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/wvuong/gogame/assets"
	"github.com/wvuong/gogame/engine"
	"golang.org/x/image/colornames"
)

type BattleScene struct {
	director *Director
	ui       *ebitenui.UI
	debug    *engine.Debug

	battlers    []*engine.Battler
	steps       []*step
	arrowSprite *engine.Sprite
	arrows      []*engine.Sprite
}

func NewBattleScene(config engine.GameConfig, state *engine.GameState, director *Director) *BattleScene {
	blackWarriorIdleSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlackWarrior_Idle, 192, 192, 8, true, frameIntervals(100*time.Millisecond, 8))
	blackWarriorRunSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlackWarrior_Run, 192, 192, 6, true, frameIntervals(100*time.Millisecond, 6))
	blackWarriorAttackSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlackWarrior_Attack, 192, 192, 4, true, frameIntervals(100*time.Millisecond, 4))
	blackWarriorGuardSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlackWarrior_Guard, 192, 192, 6, true, frameIntervals(100*time.Millisecond, 6))
	blackWarriorSprite := engine.NewSprite2(blackWarriorIdleSpriteIndex.NextFrame(), false, 0)

	blackWarriorAnimations := &engine.BattleAnimations{
		Idle:      blackWarriorIdleSpriteIndex,
		Run:       blackWarriorRunSpriteIndex,
		AttackOne: blackWarriorAttackSpriteIndex,
		Guard:     blackWarriorGuardSpriteIndex,
	}

	blueWarriorIdleSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlueWarrior_Idle, 192, 192, 8, true, frameIntervals(100*time.Millisecond, 8))
	blueWarriorRunSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlueWarrior_Run, 192, 192, 6, true, frameIntervals(100*time.Millisecond, 6))
	blueWarriorAttackSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlueWarrior_Attack, 192, 192, 4, true, frameIntervals(100*time.Millisecond, 4))
	blueWarriorGuardSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlueWarrior_Guard, 192, 192, 6, true, frameIntervals(100*time.Millisecond, 6))
	blueWarriorSprite := engine.NewSprite2(blueWarriorIdleSpriteIndex.NextFrame(), true, 0)

	blueWarriorAnimations := &engine.BattleAnimations{
		Idle:      blueWarriorIdleSpriteIndex,
		Run:       blueWarriorRunSpriteIndex,
		AttackOne: blueWarriorAttackSpriteIndex,
		Guard:     blueWarriorGuardSpriteIndex,
	}

	blackArcherIdleSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlackArcher_Idle, 192, 192, 6, true, frameIntervals(100*time.Millisecond, 6))
	blackArcherRunSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlackArcher_Run, 192, 192, 4, true, frameIntervals(100*time.Millisecond, 4))
	blackArcherAttackSpriteIndex := engine.NewHorizontalSpriteIndex(assets.BlackArcher_Attack, 192, 192, 8, false, frameIntervals(150*time.Millisecond, 8))
	blackArcherSprite := engine.NewSprite2(blackArcherIdleSpriteIndex.NextFrame(), false, 0)

	blackArcherAnimations := &engine.BattleAnimations{
		Idle:      blackArcherIdleSpriteIndex,
		Run:       blackArcherRunSpriteIndex,
		AttackOne: blackArcherAttackSpriteIndex,
	}

	root := buildUI(config)

	return &BattleScene{
		director: director,
		ui:       &ebitenui.UI{Container: root},
		debug:    &engine.Debug{Enabled: false},
		battlers: []*engine.Battler{
			engine.NewBattler(blackWarriorSprite, blackWarriorAnimations, engine.NewVector(100, 100)),
			engine.NewBattler(blueWarriorSprite, blueWarriorAnimations, engine.NewVector(400, 100)),
			engine.NewBattler(blackArcherSprite, blackArcherAnimations, engine.NewVector(100, 200)),
		},
		arrows: make([]*engine.Sprite, 0),
	}
}

func frameIntervals(duration time.Duration, count int) []time.Duration {
	intervals := make([]time.Duration, count)
	for i := range intervals {
		intervals[i] = duration
	}

	return intervals
}

func buildUI(config engine.GameConfig) *widget.Container {
	down := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(colornames.Steelblue),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:  widget.AnchorLayoutPositionEnd,
				StretchHorizontal: true,
			}),
			widget.WidgetOpts.MinSize(50, config.ScreenHeight/4),
		),
	)

	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	root.AddChild(down)

	return root
}

type step struct {
	timer  *engine.Timer
	action func()
}

func (bs *BattleScene) ChangeState(state engine.BattlerState) {
	left1 := bs.battlers[0]
	left1.SetState(state)

	left2 := bs.battlers[2]
	left2.SetState(state)

	right := bs.battlers[1]
	right.SetState(state)
}

func (bs *BattleScene) Attack() {
	actions := make([]*step, 0)
	left := bs.battlers[0]
	right := bs.battlers[1]
	startX := left.Position().X
	targetX := right.Position().X - float64(right.Dx()/3)
	distance := targetX - startX
	forwardRunningSpeed := 8.0
	backwardRunningSpeed := 10.0
	ticksNeeded := int(distance / forwardRunningSpeed)

	// move left battler towards right battler
	for range ticksNeeded {
		action := &step{
			action: func() {
				left.SetState(engine.BattlerRunning)
				left.MoveX(forwardRunningSpeed)
			},
		}

		actions = append(actions, action)
	}

	// snap left battler to target position
	actions = append(actions, &step{
		action: func() {
			left.Position().X = targetX
		},
	})

	// play attack animation and guard animation
	actions = append(actions, &step{
		action: func() {
			left.SetState(engine.BattlerAttacking)
			right.SetState(engine.BattlerGuarding)
			right.Attacked(true)
		},
	})

	// reset battlers to idle and show damage message after a delay
	actions = append(actions, &step{
		timer: engine.NewTimer(1000 * time.Millisecond),
		action: func() {
			left.SetState(engine.BattlerIdle)
			right.Attacked(false)
			right.SetState(engine.BattlerIdle)
			right.SetBottomMessage("-9999", 2*time.Second)
		},
	})

	// move left battler back to start position
	ticksNeeded = int(distance / backwardRunningSpeed)
	for range ticksNeeded {
		action := &step{
			action: func() {
				left.UseForwardAnimation(false)
				left.SetState(engine.BattlerRunning)
				left.MoveX(-backwardRunningSpeed)
			},
		}

		actions = append(actions, action)
	}

	// snap left battler to start position
	actions = append(actions, &step{
		action: func() {
			left.UseForwardAnimation(true)
			left.SetState(engine.BattlerIdle)
			left.Position().X = startX
		},
	})

	bs.steps = actions
}

func (bs *BattleScene) FireArrow() {
	actions := make([]*step, 0)
	archer := bs.battlers[2]
	startX := archer.Position().X
	startY := archer.Position().Y
	arrow := engine.NewSprite(assets.Arrow)
	arrow.ScreenPosition.X = archer.Position().X
	arrow.ScreenPosition.Y = archer.Position().Y
	arrowSpeed := 8
	distance := float64(bs.director.config.ScreenWidth) - startX
	ticks := int(distance / float64(arrowSpeed))

	// play attack animation
	actions = append(actions, &step{
		action: func() {
			archer.SetState(engine.BattlerAttacking)
		},
	})

	// hold attack animation, create arrow for rendering, change to idle
	actions = append(actions, &step{
		timer: engine.NewTimer(1000 * time.Millisecond),
		action: func() {
			bs.arrows = append(bs.arrows, arrow)

			archer.SetState(engine.BattlerIdle)
		},
	})

	// move the arrow
	for range ticks {
		actions = append(actions, &step{
			action: func() {
				// move the arrow
				x := arrow.ScreenPosition.X + float64(arrowSpeed)
				arrow.ScreenPosition.X = x
				deltaX := arrow.ScreenPosition.X - startX
				pct := deltaX / distance
				rad := pct * 3.14
				deltaY := math.Sin(rad)
				y := startY - (deltaY * 192)
				arrow.ScreenPosition.Y = y

				// figure out the position of the arrow in the next tick
				nextX := x + float64(arrowSpeed)
				nextDeltaX := deltaX + float64(arrowSpeed)
				nextPct := nextDeltaX / distance
				nextRad := nextPct * 3.14
				nextDeltaY := math.Sin(nextRad)
				nextY := startY - (nextDeltaY * 192)

				// create a triangle between current position and next position
				a := nextY - y
				b := nextX - x
				c := math.Sqrt(a*a + b*b)

				// find the angle of the triangle
				theta := math.Asin(a / c)

				// convert radians to degrees
				degrees := theta * 180 / 3.14

				// assign the arrow rotation
				arrow.Rotation = degrees
			},
		})
	}

	// delete the arrow
	actions = append(actions, &step{
		action: func() {
			idx := slices.Index(bs.arrows, arrow)
			bs.arrows = slices.Delete(bs.arrows, idx, idx+1)
		},
	})

	bs.steps = actions
}

func (bs *BattleScene) Update() {
	// go back to title screen if escape is pressed
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		bs.director.SwitchToTitle()
	}

	// toggle debug mode
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		bs.debug.Enabled = !bs.debug.Enabled
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		bs.ChangeState(engine.BattlerIdle)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		bs.ChangeState(engine.BattlerRunning)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		bs.ChangeState(engine.BattlerAttacking)
	}

	if inpututil.IsKeyJustPressed(ebiten.Key4) {
		bs.ChangeState(engine.BattlerGuarding)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		bs.Attack()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		bs.FireArrow()
	}

	// we execute one step per tick
	// execute the next pending step
	for _, step := range bs.steps {
		if step.timer == nil {
			step.action()
			bs.steps = bs.steps[1:]
			break

		} else {
			// if the step has a timer, we need to wait until the timer is ready before executing the action
			step.timer.Update()
			if step.timer.IsReady() {
				step.action()
				bs.steps = bs.steps[1:]
				break

			} else {
				break
			}
		}
	}

	// tick battlers
	for _, battler := range bs.battlers {
		battler.Update()
	}

	// tick the ui
	bs.ui.Update()
}

func (bs *BattleScene) Draw(screen *ebiten.Image) {
	// draw background
	screen.Fill(colornames.Green)

	// draw battler sprites
	for _, battler := range bs.battlers {
		battler.Draw(screen, bs.debug)
	}

	// draw arrows
	for _, arrow := range bs.arrows {
		op := &ebiten.DrawImageOptions{}
		arrow.Draw(screen, op, bs.debug)
	}

	// draw ui
	bs.ui.Draw(screen)
}

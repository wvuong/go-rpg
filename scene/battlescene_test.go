package scene

import (
	"testing"
	"time"

	"github.com/wvuong/gogame/engine"
)

// Attack and FireArrow used to assign bs.steps unconditionally, discarding an
// in-flight sequence along with the trailing steps that restore state.
//
// Both scenes here are deliberately missing their battlers and director: the
// guard has to run before either is touched, so if it is removed these panic
// rather than quietly regressing.
func TestAttackIgnoredWhileSequenceRunning(t *testing.T) {
	inFlight := []*step{{action: func() {}}, {action: func() {}}}
	bs := &BattleScene{steps: inFlight}

	bs.Attack()

	if len(bs.steps) != len(inFlight) {
		t.Fatalf("in-flight queue length = %d, want %d (Attack replaced a running sequence)", len(bs.steps), len(inFlight))
	}
}

func TestFireArrowIgnoredWhileSequenceRunning(t *testing.T) {
	inFlight := []*step{{action: func() {}}, {action: func() {}}}
	bs := &BattleScene{steps: inFlight}

	bs.FireArrow()

	if len(bs.steps) != len(inFlight) {
		t.Fatalf("in-flight queue length = %d, want %d (FireArrow replaced a running sequence)", len(bs.steps), len(inFlight))
	}
}

func TestAdvanceStepsRunsOneStepPerTick(t *testing.T) {
	var order []int
	bs := &BattleScene{}
	for i := range 3 {
		n := i
		bs.steps = append(bs.steps, &step{action: func() { order = append(order, n) }})
	}

	for tick := 1; tick <= 3; tick++ {
		bs.advanceSteps()
		if len(order) != tick {
			t.Fatalf("after %d tick(s) ran %d step(s), want %d", tick, len(order), tick)
		}
	}

	if bs.busy() {
		t.Error("queue still busy after draining every step")
	}

	for i, got := range order {
		if got != i {
			t.Errorf("step %d ran at position %d; steps must run in order", got, i)
		}
	}
}

// A timed step holds the queue until its timer is ready.
func TestAdvanceStepsBlocksOnTimer(t *testing.T) {
	ran := false
	bs := &BattleScene{steps: []*step{
		{timer: engine.NewTimer(100 * time.Millisecond), action: func() { ran = true }},
	}}

	// engine.Timer converts the duration to ticks via ebiten.TPS() (60), so
	// 100ms is 6 ticks.
	for range 5 {
		bs.advanceSteps()
		if ran {
			t.Fatal("timed step ran before its timer was ready")
		}
	}

	bs.advanceSteps()
	if !ran {
		t.Error("timed step never ran once its timer was ready")
	}
	if bs.busy() {
		t.Error("completed step was not removed from the queue")
	}
}

func TestAdvanceStepsOnEmptyQueue(t *testing.T) {
	bs := &BattleScene{}
	bs.advanceSteps() // must not panic
	if bs.busy() {
		t.Error("empty queue reported as busy")
	}
}

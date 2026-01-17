package game

type Ticker struct {
	count int
}

func NewTicker() *Ticker {
	return &Ticker{}
}

func (t *Ticker) Increment() {
	t.count++
}

func (t *Ticker) Count() int {
	return t.count
}

func (t *Ticker) Reset() {
	t.count = 0
}

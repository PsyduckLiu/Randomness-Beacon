package consensus

import "time"

// timer
const StateTimerOut = 5 * time.Second

var timeArray []float64
var lastTime time.Time
var outputNum int = 0

type RequestTimer struct {
	*time.Ticker
	IsOk bool
}

// initialize timer
func newRequestTimer() *RequestTimer {
	tick := time.NewTicker(StateTimerOut)
	tick.Stop()
	return &RequestTimer{
		Ticker: tick,
		IsOk:   false,
	}
}

// start a timer
func (rt *RequestTimer) tick(time time.Duration) {
	if rt.IsOk {
		return
	}
	rt.Reset(time)
	rt.IsOk = true
}

// stop a timer
func (rt *RequestTimer) tack() {
	rt.IsOk = false
	rt.Stop()
}

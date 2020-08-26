package main

import (
	"sync"
	"time"
)

// Alarm will tick for a duration or until ended by another process.
type Alarm struct {
	mtx         sync.Mutex
	ticker      *time.Ticker
	currentTick float64
	duration    time.Duration
	stopped     chan bool
	ring        chan bool
}

// NewAlarm creates a new alarm but doesn't start the ticking
func NewAlarm(duration time.Duration) *Alarm {
	alarm := &Alarm{
		ticker: time.NewTicker(time.Second),
		duration: duration,
		ring: make(chan bool),
	}

	return alarm
}

// Start the ticking the clock
func (a *Alarm) Start() {
	// Begin ticking
	go func() {
		for {
			select {
			case <-a.stopped:
				return
			case <-a.ticker.C:
				a.mtx.Lock()
				a.currentTick++
				a.mtx.Unlock()
			default:
				if a.currentTick >= a.duration.Seconds() {
					a.ring <- true
					a.Stop()
				}
			}
		}
	}()
}

// Stop the alarm
func (a *Alarm) Stop() {
	a.ticker.Stop()
	a.stopped <- true
}

// Restart the alarm
func (a *Alarm) Reset() {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.currentTick = 0
	a.ticker = time.NewTicker(time.Second)
}


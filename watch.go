package main

import "time"

type Watcher struct {
	*Page
	tick  <-chan time.Time
	walkQ chan *Page
}

var (
	countWatchers int
)

// NewWatcher will schedule the page for periodic walks starting
// immediately.
func NewWatcher(p *Page) (w *Watcher) {
	countWatchers++
	return &Watcher{
		Page: p,
		tick: time.Tick(1 * time.Minute),
	}
}

func (w *Watcher) StartWatching() {
	for now := range w.tick {

		// Send the page to the walk queue, the channel will handle
		// the proper queuing queue
		w.Page.TimeStamp = Timestamp(now)
		walkQ <- w.Page
	}
}

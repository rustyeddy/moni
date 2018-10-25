package moni

type Scheduler struct {
	URLQ   chan string
	CrawlQ chan *Page
	ErrorQ chan error
}

var (
	sched *Scheduler
)

func GetScheduler() *Scheduler {
	if sched == nil {
		sched = &Scheduler{
			URLQ:   make(chan string),
			CrawlQ: make(chan *Page),
			ErrorQ: make(chan error),
		}
	}
	return sched
}

// Start the scheduler
func (sched *Scheduler) Start() {
	go urlWatcher(sched.URLQ, sched.CrawlQ, sched.ErrorQ)
	go crawlWatcher(sched.CrawlQ, sched.ErrorQ)
	go errorWatcher(sched.ErrorQ)
}

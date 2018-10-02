package inv

import "time"

type PageInfo struct {
	URL        string
	StatusCode int
	Links      map[string]int
	Start      time.Time
	End        time.Time
}

package moni

import (
	"time"
)

const (
	hidNone = iota
	hidErrorWatch
)

type Event struct {
	id  string // message comes from some where
	msg string
	num int64
	time.Time
}

type History []Event

var history History

func RecordEvent(id string, msg string, num int64) History {
	history = append(history, Event{
		id:   id,
		msg:  msg,
		num:  num,
		Time: time.Now(),
	})
	return history
}
func RecordMesg(msg string) History { return RecordEvent(msg, msg, 0) }
func RecordNum(num int64) History   { return RecordEvent("FOO", "", 0) }

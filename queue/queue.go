package queue

import (
    "fmt"
)

type QueuedCmd func() error

type Queue struct {
	queue []QueuedCmd
}



func NewQueue() *Queue {
	q := &Queue{}
	q.queue = make([]QueuedCmd, 0, 4)
	return q
}

func (q *Queue) Add(qcmds ...QueuedCmd) {
	for _, qcmd := range qcmds {
		q.queue = append(q.queue, qcmd)
	}
}

func (q *Queue) Run() (err error) {
    defer func() {
        if e := recover(); e != nil {
            err = fmt.Errorf("panicked: %v", e)
        }
    }()
	for _, qcmd := range q.queue {
        if err = qcmd(); err != nil {
            return
		}
	}
	return
}

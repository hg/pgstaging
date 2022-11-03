package worker

import (
	"github.com/hg/pgstaging/worker/command"
	"fmt"
	"time"
)

type Action int

const (
	ActionStart = iota
	ActionStop
	ActionCreate
	ActionDrop
)

type task struct {
	do   Action
	name string
	done chan<- error
}

type Worker struct {
	tasks chan *task
}

func New() *Worker {
	w := &Worker{
		tasks: make(chan *task, 100),
	}
	go w.handle()
	return w
}

func (w *Worker) Enqueue(act Action, name string) <-chan error {
	done := make(chan error, 1)
	w.tasks <- &task{act, name, done}
	return done
}

func (w *Worker) handle() {
	for ts := range w.tasks {

		// todo: remove
		time.Sleep(2 * time.Second)

		switch ts.do {
		case ActionStart:
			ts.done <- command.StartCluster(ts.name)
		case ActionStop:
			ts.done <- command.StopCluster(ts.name)
		case ActionCreate:
			ts.done <- command.CreateCluster(ts.name)
		case ActionDrop:
			ts.done <- command.DropCluster(ts.name)
		default:
			ts.done <- fmt.Errorf("invalid action %v", ts.do)
		}
	}
}

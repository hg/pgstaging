package worker

import (
	"github.com/hg/pgstaging/config"
	"github.com/hg/pgstaging/worker/command"
	"fmt"
	"log"
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
	pass string
	done chan<- error
}

type Worker struct {
	tasks chan *task
	conf  *config.Config
}

func New(conf *config.Config) *Worker {
	w := &Worker{
		tasks: make(chan *task, 100),
		conf:  conf,
	}
	go w.handle()
	return w
}

func (w *Worker) Enqueue(act Action, name, pass string) <-chan error {
	done := make(chan error, 1)
	w.tasks <- &task{act, name, pass, done}
	return done
}

func (w *Worker) handle() {
	for ts := range w.tasks {
		log.Printf("starting task [%d] on %s", ts.do, ts.name)

		switch ts.do {
		case ActionCreate:
			ts.done <- command.CreateCluster(ts.name, ts.pass)

		case ActionStart, ActionStop, ActionDrop:
			handleModification(ts, w.conf.Passwd)

		default:
			ts.done <- fmt.Errorf("invalid action %v", ts.do)
		}
	}
}

func handleModification(ts *task, admin string) {
	if !validPassword(ts.name, ts.pass, admin) {
		ts.done <- fmt.Errorf("invalid password")
		return
	}

	switch ts.do {
	case ActionStart:
		ts.done <- command.StartCluster(ts.name)
	case ActionStop:
		ts.done <- command.StopCluster(ts.name)
	case ActionDrop:
		ts.done <- command.DropCluster(ts.name)
	default:
		panic(fmt.Sprintf("bad action %v", ts.do))
	}
}

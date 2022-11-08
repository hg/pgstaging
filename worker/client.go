package worker

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"sync/atomic"
)

var taskId = atomic.Int32{}

type Client struct {
	tasks chan *task
	dones chan chan<- error
	in    io.Reader
	out   io.Writer
}

func NewClient(in io.Reader, out io.Writer) *Client {
	wrk := &Client{
		tasks: make(chan *task, 100),
		dones: make(chan chan<- error, 100),
		in:    in,
		out:   out,
	}
	go wrk.sendTasks()
	go wrk.receiveResults()
	return wrk
}

func (w *Client) sendTasks() {
	for ts := range w.tasks {
		js, err := json.Marshal(&ts.pay)
		if err != nil {
			log.Fatalf("error marshalling task: %v", err)
		}
		if _, err = w.out.Write(js); err == nil {
			_, err = w.out.Write([]byte("\n"))
		}
		if err != nil {
			log.Fatalf("error writing task: %v", err)
		}
	}
}

func (w *Client) receiveResults() {
	read := bufio.NewReader(w.in)

	for nextId := int32(1); ; nextId++ {
		line, err := read.ReadBytes('\n')
		if err != nil {
			log.Fatalf("error reading task result: %v", err)
		}

		var res result
		if err = json.Unmarshal(line, &res); err != nil {
			log.Fatalf("error parsing task result: %v", err)
		}

		if res.Id != nextId {
			log.Fatalf("unexpected task id %d (want %d)", res.Id, nextId)
		}

		done := <-w.dones

		if res.Error == "" {
			done <- nil
		} else {
			done <- errors.New(res.Error)
		}
	}
}

func (w *Client) Enqueue(act Action, name, pass string) <-chan error {
	done := make(chan error, 1)
	w.tasks <- &task{
		pay: payload{
			Id:   taskId.Add(1),
			Do:   act,
			Name: name,
			Pass: pass,
		},
		done: done,
	}
	w.dones <- done
	return done
}

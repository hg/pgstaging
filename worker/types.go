package worker

import "github.com/hg/pgstaging/worker/command"

type Action int

const (
	ActionStart = iota
	ActionStop
	ActionCreate
	ActionForceCreate
	ActionDrop
)

type payload struct {
	Id   int32  `json:"id"`
	Do   Action `json:"do"`
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type task struct {
	pay  payload
	done chan<- command.Result
}

type result struct {
	Id    int32  `json:"id"`
	Error string `json:"error"`
	Data  any    `json:"data"`
}

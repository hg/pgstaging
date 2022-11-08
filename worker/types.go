package worker

type Action int

const (
	ActionStart = iota
	ActionStop
	ActionCreate
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
	done chan<- error
}

type result struct {
	Id    int32  `json:"id"`
	Error string `json:"error"`
}

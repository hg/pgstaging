package worker

import (
	"bufio"
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker/command"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Server struct {
	pass string
	read *bufio.Reader
}

func NewServer(pass string) *Server {
	return &Server{
		pass: pass,
		read: bufio.NewReader(os.Stdin),
	}
}

func (srv *Server) Run() error {
	for {
		pay, err := srv.recvCommand()
		log.Printf("received command %v", pay)

		if err = srv.processTask(pay); err == nil {
			srv.sendResult(result{pay.Id, ""})
		} else {
			srv.sendResult(result{pay.Id, err.Error()})
		}
	}
}

func (srv *Server) recvCommand() (*payload, error) {
	line, err := srv.read.ReadBytes('\n')
	if err != nil {
		log.Fatalf("error reading command: %v", err)
	}

	var pay payload

	if err = json.Unmarshal(line, &pay); err != nil {
		log.Fatalf("error parsing command: %v", err)
	}

	// In case the main application was broken and is used to
	// pass something like /etc/passwd or ../../../../
	if !util.IsDevName(pay.Name) {
		log.Fatalf("bad name received")
	}
	if pay.Pass != "" && !util.IsOkPassword(pay.Pass) {
		log.Fatalf("bad password received")
	}

	return &pay, err
}

func (srv *Server) sendResult(res result) {
	js, err := json.Marshal(&res)
	if err != nil {
		log.Fatalf("error serializing task result: %v", err)
	}

	if _, err = os.Stdout.Write(js); err == nil {
		_, err = os.Stdout.Write([]byte("\n"))
	}

	if err != nil {
		log.Fatalf("error writing task result: %v", err)
	}
}

func (srv *Server) processTask(ts *payload) error {
	log.Printf("starting task [%d] on %s", ts.Do, ts.Name)

	switch ts.Do {
	case ActionCreate:
		return command.CreateCluster(ts.Name, ts.Pass)

	case ActionStart, ActionStop, ActionDrop:
		return handleModification(ts, srv.pass)

	default:
		return fmt.Errorf("invalid action %v", ts.Do)
	}
}

func handleModification(ts *payload, admin string) error {
	if !validPassword(ts.Name, ts.Pass, admin) {
		return fmt.Errorf("invalid password")
	}

	switch ts.Do {
	case ActionStart:
		return command.StartCluster(ts.Name)
	case ActionStop:
		return command.StopCluster(ts.Name)
	case ActionDrop:
		return command.DropCluster(ts.Name)
	default:
		panic(fmt.Sprintf("bad action %v", ts.Do))
	}
}

package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/hg/pgstaging/web/util"
	"github.com/hg/pgstaging/worker/command"
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
		pay := srv.recvCommand()
		log.Printf("received command %v", pay)

		if re := srv.processTask(pay); re.Err == nil {
			srv.sendResult(result{pay.Id, "", re.Data})
		} else {
			srv.sendResult(result{pay.Id, re.Err.Error(), nil})
		}
	}
}

func (srv *Server) recvCommand() *payload {
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

	return &pay
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

func (srv *Server) processTask(ts *payload) command.Result {
	log.Printf("starting task [%d] on %s", ts.Do, ts.Name)

	switch ts.Do {
	case ActionCreate:
		return command.CreateCluster(ts.Name, ts.Pass, false)

	case ActionForceCreate:
		return command.CreateCluster(ts.Name, ts.Pass, true)

	case ActionStart, ActionStop, ActionDrop:
		return handleModification(ts, srv.pass)

	default:
		return command.Result{Err: fmt.Errorf("invalid action %v", ts.Do)}
	}
}

func handleModification(ts *payload, admin string) command.Result {
	if !validPassword(ts.Name, ts.Pass, admin) {
		return command.Result{Err: fmt.Errorf("invalid password")}
	}

	switch ts.Do {
	case ActionStart:
		return command.Result{Err: command.StartCluster(ts.Name)}
	case ActionStop:
		return command.Result{Err: command.StopCluster(ts.Name)}
	case ActionDrop:
		return command.Result{Err: command.DropCluster(ts.Name)}
	default:
		panic(fmt.Sprintf("bad action %v", ts.Do))
	}
}

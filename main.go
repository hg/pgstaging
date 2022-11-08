package main

import (
	"github.com/hg/pgstaging/config"
	"github.com/hg/pgstaging/service"
	"github.com/hg/pgstaging/util"
	"github.com/hg/pgstaging/web"
	"github.com/hg/pgstaging/worker"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

//go:embed templates/* assets/*
var files embed.FS

func sanityChecks() error {
	if runtime.GOOS != "linux" {
		return errors.New("unsupported operating system")
	}
	if os.Geteuid() != 0 {
		return errors.New("server needs to run as root")
	}
	ls := exec.Command("pg_lsclusters", "--help")
	return ls.Run()
}

func showUsage() {
	bin, _ := os.Executable()

	fmt.Printf(`
usage: %s install|run
	install: install binary and system service
	run: start the server
`, bin)

	os.Exit(1)
}

func startServer() error {
	if err := sanityChecks(); err != nil {
		log.Fatalf("init failed: %v", err)
	}

	conf, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	wrk, err := startWorker(conf)
	if err != nil {
		return err
	}

	uid, err := util.GetUserId(conf.User)
	if err != nil {
		return err
	}
	err = syscall.Seteuid(int(uid.UID))
	if err != nil {
		return err
	}

	return web.Start(conf.Listen, wrk, files)
}

func startWorker(conf *config.Config) (*worker.Client, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(exe, "worker", conf.Passwd)

	inR, inW, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	outR, outW, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	cmd.Stdin = outR
	cmd.Stdout = inW
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start worker: %v", err)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Fatalf("worker exited with error: %v", err)
		}
	}()

	return worker.NewClient(inR, outW), nil
}

func runWorker() error {
	passwd := os.Args[len(os.Args)-1]
	if passwd == "" {
		log.Fatalf("admin password not specified (invalid call?)")
	}
	srv := worker.NewServer(passwd)
	return srv.Run()
}

func doCommand(name string) (err error) {
	switch name {
	case "install":
		if err = service.Install(); err == nil {
			err = service.Enable()
		}

	case "run":
		err = startServer()

	case "worker":
		err = runWorker()

	default:
		showUsage()
	}

	return
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
	}
	if err := doCommand(os.Args[1]); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}

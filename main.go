package main

import (
	"github.com/hg/pgstaging/service"
	"github.com/hg/pgstaging/web"
	"github.com/hg/pgstaging/worker"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
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
	wrk := worker.New()
	return web.Start(":80", wrk, files)
}

func doCommand(name string) (err error) {
	switch name {
	case "install":
		if err = service.Install(); err == nil {
			err = service.Enable()
		}

	case "run":
		err = startServer()

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

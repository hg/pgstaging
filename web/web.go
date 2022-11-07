package web

import (
	"github.com/hg/pgstaging/web/sessions"
	"github.com/hg/pgstaging/worker"
	"embed"
	"html/template"
	"log"
	"net/http"
)

func loadTemplate(files embed.FS) *template.Template {
	tpl, err := template.ParseFS(files, "templates/index.tmpl")
	if err != nil {
		log.Fatal("could not load template", err)
	}
	return tpl
}

func Start(address string, wrk *worker.Worker, files embed.FS) error {
	srv := &server{
		sessions: sessions.New(),
		worker:   wrk,
		tpl:      loadTemplate(files),
	}

	mux := http.NewServeMux()

	// static files
	mux.Handle("/assets/", http.FileServer(http.FS(files)))

	mux.HandleFunc("/api/create", srv.wrap(serveCreate))
	mux.HandleFunc("/api/modify", srv.wrap(serveModify))
	mux.HandleFunc("/", srv.wrap(serveIndex))

	log.Printf("listening on %s", address)
	return http.ListenAndServe(address, mux)
}

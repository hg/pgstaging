package main

import (
	"github.com/hg/pgstaging/pg"
	"github.com/hg/pgstaging/util"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

//go:embed templates/* assets/*
var files embed.FS

func loadTemplate() *template.Template {
	tpl, err := template.ParseFS(files, "templates/index.tmpl")
	if err != nil {
		log.Fatal("could not load template", err)
	}
	return tpl
}

var tpl = loadTemplate()
var reNonAlnum = regexp.MustCompile(`[^a-z0-9_]]`)

var muSessions = sync.Mutex{}
var sessions = make(map[string]*session)

type pageModel struct {
	Result   string
	Message  string
	Clusters []clusterModel
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/", "/index.html":
		break
	default:
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	model := pageModel{
		Clusters: clustersToViewModels(pg.GetActiveClusters()),
	}

	cookie, err := r.Cookie("session")
	if err == nil {
		muSessions.Lock()
		sess := sessions[cookie.Value]
		muSessions.Unlock()

		if sess != nil {
			switch sess.status {
			case statusNone:
			case statusQueued:
				model.Result = "Queued"
			case statusError:
				model.Result = "Error"
			}
			model.Message = sess.message
		}
	}

	w.Header().Set("Content-Type", "text/html")

	err = tpl.Execute(w, model)
	if err != nil {
		log.Print("error writing response", err)
	}
}

type clusterModel struct {
	Name    string
	Port    uint16
	User    string
	Pass    string
	Dev     bool
	Running bool
}

func clustersToViewModels(clusters []pg.Cluster) (result []clusterModel) {
	for _, cluster := range clusters {
		result = append(result, clusterModel{
			Name:    cluster.Cluster,
			Port:    cluster.Port,
			User:    "sc",
			Pass:    "sc",
			Dev:     strings.HasPrefix(cluster.Cluster, "dev_"),
			Running: cluster.Running != 0,
		})
	}
	return
}

type status int

const (
	statusNone = iota
	statusQueued
	statusError
)

type session struct {
	status  status
	message string
}

func setResult(sessionId string, stat status, message string) {
	muSessions.Lock()
	sessions[sessionId] = &session{
		status:  stat,
		message: message,
	}
	muSessions.Unlock()
}

func getOrCreateSession(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}
	cookie = &http.Cookie{
		Name:     "session",
		Value:    util.RandomString(16),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	return cookie.Value
}

type task struct{}

var taskQueue = make(chan *task, 100)

func createCluster(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	name = strings.ToLower(name)
	name = reNonAlnum.ReplaceAllString(name, "")

	sessionId := getOrCreateSession(w, r)

	if name == "" || len(name) > 32 {
		setResult(sessionId, statusError, fmt.Sprintf("некорректное имя '%s'", name))
	}

	taskQueue <- &task{}
	setResult(sessionId, statusQueued, "")

	http.Redirect(w, r, "/", http.StatusFound)
}

func startWebServer() {
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.FileServer(http.FS(files)))
	mux.HandleFunc("/api/create", createCluster)
	mux.HandleFunc("/", serveIndex)

	srv := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	err := srv.ListenAndServe()
	log.Fatal(err)
}

func main() {
	startWebServer()
}

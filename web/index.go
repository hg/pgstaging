package web

import (
	"github.com/hg/pgstaging/pg"
	"github.com/hg/pgstaging/web/sessions"
	"github.com/hg/pgstaging/web/util"
	"log"
	"net/http"
	"os"
	"time"
)

type clusterModel struct {
	Name     string
	Port     uint16
	User     string
	Pass     string
	Dev      bool
	Running  bool
	Modified time.Time
}

type event struct {
	When    time.Time
	Error   bool
	Status  string
	Message string
}

type pageModel struct {
	Clusters []clusterModel
	Events   []event
}

func lastModified(name string) time.Time {
	st, err := os.Stat(name)
	if err != nil {
		return time.Time{}
	}
	return st.ModTime()
}

func clustersToViewModels(clusters []pg.Cluster) (result []clusterModel) {
	for _, cluster := range clusters {
		dev := util.IsDevName(cluster.Cluster)
		mdl := clusterModel{
			Name:     cluster.Cluster,
			Port:     cluster.Port,
			Dev:      dev,
			Running:  cluster.Running != 0,
			Modified: lastModified(cluster.ConfigDir),
		}
		if dev {
			mdl.User = "sc"
			mdl.Pass = "sc"
		}
		result = append(result, mdl)
	}
	return
}

func humanizeStatus(status string) string {
	switch status {
	case "queued":
		return "Задача поставлена в очередь"
	case "success":
		return "Задача выполнена успешно"
	case "error":
		return "Не удалось выполнить задачу"
	default:
		return status
	}
}

func eventsToViewModel(events []sessions.Event) []event {
	var out []event
	for _, evt := range events {
		out = append(out, event{
			When:    evt.Created,
			Error:   evt.Status == "error",
			Status:  humanizeStatus(evt.Status),
			Message: evt.Message,
		})
	}
	return out
}

func serveIndex(rc *requestContext) {
	switch rc.request.URL.Path {
	case "/", "/index.html":
		break
	default:
		http.NotFound(rc.writer, rc.request)
		return
	}

	if !rc.isMethod(http.MethodGet) {
		return
	}

	model := pageModel{
		Clusters: clustersToViewModels(pg.GetActiveClusters()),
		Events:   eventsToViewModel(rc.session.Events()),
	}

	rc.writer.Header().Set("Content-Type", "text/html")
	err := rc.srv.tpl.Execute(rc.writer, model)

	if err != nil {
		log.Print("error writing response", err)
	}
}

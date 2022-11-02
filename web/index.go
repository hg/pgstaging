package web

import (
	"github.com/hg/pgstaging/pg"
	"log"
	"net/http"
	"strings"
)

type clusterModel struct {
	Name    string
	Port    uint16
	User    string
	Pass    string
	Dev     bool
	Running bool
}

type pageModel struct {
	Status   string
	Message  string
	Clusters []clusterModel
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
		Status:   (rc.session.Get("status")).(string),
		Message:  (rc.session.Get("message")).(string),
	}

	rc.writer.Header().Set("Content-Type", "text/html")
	err := rc.srv.tpl.Execute(rc.writer, model)

	if err != nil {
		log.Print("error writing response", err)
	}
}

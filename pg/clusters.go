package pg

import (
	"encoding/json"
	"log"
	"os/exec"
	"sort"
)

type Cluster struct {
	Version   string `json:"version"`
	Running   int    `json:"running"`
	Cluster   string `json:"cluster"`
	Port      uint16 `json:"port,string"`
	ConfigDir string `json:"configdir"`
	PgData    string `json:"pgdata"`
}

func GetActiveClusters() (result []Cluster) {
	cmd := exec.Command("ssh", "sc_pg", "pg_lsclusters", "--json")

	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("could not read lsclusters output: %v", err)
		return
	}
	err = json.Unmarshal(stdout, &result)
	if err != nil {
		log.Printf("could not parse lsclusters output: %v", err)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Port < result[j].Port
	})
	return
}

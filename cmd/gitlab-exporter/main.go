package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	token = "___"
	gitlabUrl = "https://git.acme.org"
)


func main() {

	// Update the stats regularly
	stats := getStats()
	go func() {
		for range time.Tick(5 * time.Minute) {
			stats = getStats()
		}
	}()

	// Start http server that responds with the stats
	router := http.NewServeMux()
	router.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", stats)
	})

	fmt.Printf("Listening on port :8123");
	if err := http.ListenAndServe(":8123", router); err != nil {
		panic(fmt.Sprintf("Could not bind webserver to %q: %v", ":8123", err))
	}
}

func getStats() string {
	fmt.Printf("Updating")
	projects := GetRepositories(gitlabUrl, token)
	stats := fmt.Sprintf("gitlab_last_update %d\n", time.Now().Unix())

	for _,project := range projects {
		stats = fmt.Sprintf("%s\n%s", stats, project.PrometheusStats())
	}
	return stats
}

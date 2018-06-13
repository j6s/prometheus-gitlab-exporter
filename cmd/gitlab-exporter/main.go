package main

import (
	"fmt"
	"net/http"
	"time"
	"flag"
)

var (
	token = flag.String(
		"token",
		"",
		"The personal access token that is used to communicate with the API." +
			"for more information about the Gitlab API please refer to the" +
			"[Gitlab documentation](https://docs.gitlab.com/ee/api/#personal-access-tokens)",
	)
	gitlabUrl = flag.String(
		"url",
		"",
		"The base URL of the gitlab instance including protocol.\n" +
			"This string must not contain any path information other than the\n" +
			"index route of gitlab. If your server runs on a non-standard port\n" +
			"(not 80 or 443 for http and https) then you may specify it using\n" +
			"a colon.\n\n" +
			"Examples:\n" +
		 	"\t- `https://git.acme.org` - gitlab running on port 443 directly" +
			"\t- `http://git.acme.org:8080` - gitlab running on port 8080 over http" +
			"\t- `https://acme.org/gitlab` - gitlab running in a subdirectory",
	)
)


func main() {
	flag.Parse();

	if *gitlabUrl == "" && *token == "" {
		panic("The arguments --token and --url must be set")
	}

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
	projects := GetRepositories(*gitlabUrl, *token)
	stats := fmt.Sprintf("gitlab_last_update %d\n", time.Now().Unix())

	for _,project := range projects {
		stats = fmt.Sprintf("%s\n%s", stats, project.PrometheusStats())
	}
	return stats
}

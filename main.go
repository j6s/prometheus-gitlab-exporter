package main

import (
	"fmt"
	"net/http"
	"time"
	"flag"
	"log"
	"os"
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
	pollInterval = flag.Duration(
		"poll-interval",
		5 * time.Minute,
		"Poll interval in minutes. The data will be updated every time interval\n" +
			"in order to avoid excessive API use.\n" +
			"Every string accepted by the " +
			"(https://golang.org/pkg/time/#example_Duration)[golang time package] is valid.\n" +
			"Default to '5m'",
	)
	bind = flag.String(
		"bind",
		":8123",
		"Address to bind to.\n" +
			"The service will be available at this address.",
	)
)


func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Prometheus Gitlab exporter \n")
		fmt.Fprintf(os.Stderr, "============================\n")
		fmt.Fprintf(os.Stderr, "Simple exporter that exposes gitlab project statistics to prometheus.\n")
		fmt.Fprintf(os.Stderr, "https://github.com/j6s/prometheus-gitlab-exporter\n")
		fmt.Fprintf(os.Stderr, "\nUsage\n")
		fmt.Fprintf(os.Stderr,   "-----\n")
		fmt.Fprintf(os.Stderr, "$ prometheus-gitlab-exporter --url='https://git.acme.org' --token='abcdef123'\n")
		fmt.Fprintf(os.Stderr, "$ prometheus-gitlab-exporter --url='https://git.acme.org' --token='abcdef123' --poll-interval='15m' --bind='hostname.com:9898'\n")
		fmt.Fprintf(os.Stderr, "\nArguments\n")
		fmt.Fprintf(os.Stderr,   "---------\n")
		flag.PrintDefaults()
	}
	flag.Parse();

	// Validate arguments
	if *gitlabUrl == "" && *token == "" {
		log.Fatalf("The arguments --token and --url must be set")
	}

	// Update the stats regularly
	stats := "";
	go func() {
		stats = getStats()
		log.Printf("Updating data every %v seconds", pollInterval.Seconds())
		for range time.Tick(*pollInterval) {
			stats = getStats()
		}
	}()

	// Start http server that responds with the stats
	router := http.NewServeMux()
	router.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s", stats)
	})

	log.Printf("Listening on port :8123");
	if err := http.ListenAndServe(*bind, router); err != nil {
		log.Printf("Could not bind webserver to %q: %v", *bind, err)
	}
}

func getStats() string {
	log.Printf("Updating")
	projects := GetRepositories(*gitlabUrl, *token)
	stats := fmt.Sprintf("gitlab_last_update %d\n", time.Now().Unix())

	for _,project := range projects {
		stats = fmt.Sprintf("%s\n%s", stats, project.PrometheusStats())
	}
	return stats
}

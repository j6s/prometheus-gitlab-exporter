package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ProjectStats struct {
	CommitCount      int `json:"commit_count"`
	StorageSize      int `json:"storage_size"`
	RepositorySize   int `json:"repository_size"`
	LfsObjectSize    int `json:"lfs_object_size"`
	JobArtifactsSize int `json:"job_artifacts_size"`
}

type Project struct {
	PathWithNamespace string       `json:"path_with_namespace"`
	StarCount         int          `json:"star_count"`
	ForkCount         int          `json:"fork_count"`
	OpenIssueCount    int          `json:"openIssueCount"`
	LastActivityAt    time.Time    `json:"last_activity_at"`
	Statistics        ProjectStats `json:"statistics"`
}

/**
 *	Extracts the stats from a single project into a
 *	prometheus compatible string.
 *
 *	@param project The project to extract stats from
 *	@return A prometheus style statistics string for the project
 */
func (project Project) PrometheusStats() string {
	stats := ""
	stats = fmt.Sprintf("%s\ngitlab_project_stars{repo=%s} %d", stats, project.PathWithNamespace, project.StarCount)
	stats = fmt.Sprintf("%s\ngitlab_project_forks{repo=%s} %d", stats, project.PathWithNamespace, project.ForkCount)
	stats = fmt.Sprintf("%s\ngitlab_project_commit_count{repo=%s} %d", stats, project.PathWithNamespace, project.Statistics.CommitCount)
	stats = fmt.Sprintf("%s\ngitlab_project_storage_size{repo=%s} %d", stats, project.PathWithNamespace, project.Statistics.StorageSize)
	stats = fmt.Sprintf("%s\ngitlab_project_repository_size{repo=%s} %d", stats, project.PathWithNamespace, project.Statistics.RepositorySize)
	stats = fmt.Sprintf("%s\ngitlab_project_lfs_object_size{repo=%s} %d", stats, project.PathWithNamespace, project.Statistics.LfsObjectSize)
	stats = fmt.Sprintf("%s\ngitlab_project_job_artifacts_size{repo=%s} %d", stats, project.PathWithNamespace, project.Statistics.JobArtifactsSize)
	stats = fmt.Sprintf("%s\ngitlab_project_last_activity{repo=%s} %d", stats, project.PathWithNamespace, project.LastActivityAt.Unix())
	return stats
}

/**
 *	Fetches all projects from the configured gitlab endpoint.
 *	This method will panic if there are any errors.
 *
 *	@return A list of all projects known to the current gitlab instance
 *	@panic If there is an error while fetching the projects or decoding the JSON response
 */
func GetRepositories(gitlabUrl string, token string) []Project {
	projects := make([]Project, 0)
	page := 1

	for true {
		// Fetch a page from the API.
		projectsUrl := fmt.Sprintf("%s/api/v4/projects?private_token=%s&per_page=100&statistics=1&page=%d", gitlabUrl, token, page)
		log.Printf("Requesting %s\n", projectsUrl)
		response, error := http.Get(projectsUrl)
		if error != nil {
			panic(error)
		}

		// Merge the results back to the complete array.
		projectsInPage := make([]Project, 0)
		error = json.NewDecoder(response.Body).Decode(&projectsInPage)
		if error != nil {
			log.Fatalf(error.Error())
		}
		projects = append(projects, projectsInPage...)

		// Parse the X-Next-Page Header in order to figure out
		// if another page should be requested.
		// If not, then return the current projects.
		pageHeader := response.Header["X-Next-Page"][0]
		if pageHeader == "" {
			return projects
		}
		page, error = strconv.Atoi(pageHeader)
		if error != nil {
			log.Fatalf(error.Error())
		}
	}

	return projects
}

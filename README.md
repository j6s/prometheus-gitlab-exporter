# Gitlab exporter

Simple prometheus exporter that exposes basic stats about
projects in gitlab.

## Installation

To install gitlab exporter simply get the package using `go get github.com/j6s/prometheus-gitlab-exporter`.

## Usage

A personal access token is required in order to communicate with the API.
You can learn more about personal access tokens and how to manage them [in gitlab API documentation](https://docs.gitlab.com/ee/api/#personal-access-tokens).

```
$ prometheus-gitlab-exporter --help
Prometheus Gitlab exporter
============================
Simple exporter that exposes gitlab project statistics to prometheus.
https://github.com/j6s/prometheus-gitlab-exporter

Usage
-----
$ prometheus-gitlab-exporter --url='https://git.acme.org' --token='abcdef123'
$ prometheus-gitlab-exporter --url='https://git.acme.org' --token='abcdef123' --poll-interval='15m' --bind='hostname.com:9898'

Arguments
---------
  -bind string
    	Address to bind to. The service will be available at this address.
    	 (default ":8123")
  -poll-interval duration
    	Poll interval in minutes. The data will be updated every time interval
    	in order to avoid excessive API use.
    	Every string accepted by the golang time package is valid.
    	https://golang.org/pkg/time/#example_Duration
    	 (default 5m0s)
  -token string
    	The personal access token that is used to communicate with
    	the API. For more information about the Gitlab API please
    	refer to the GitLab documentation.
    	https://docs.gitlab.com/ee/api/#personal-access-tokens
  -url string
    	The base URL of the gitlab instance including protocol.
    	This string must not contain any path information other than the
    	index route of gitlab. If your server runs on a non-standard port
    	(not 80 or 443 for http and https) then you may specify it using
    	a colon.
```

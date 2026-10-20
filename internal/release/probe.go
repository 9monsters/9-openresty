package release

import (
	"net/http"
	"regexp"
	"strings"
	"time"
)

// serverNameRE extracts server_name values from nginx config content.
var serverNameRE = regexp.MustCompile(`server_names+([^;]+);`)

// healthCheck sends an HTTP GET to each server_name found in the changed files.
// It returns true when all probes return HTTP status < 500.
func (c *Coordinator) healthCheck(files []ReleaseFile) bool {
	names := extractServerNames(files)
	if len(names) == 0 {
		// No server_name directives changed — probe localhost as a basic liveness check.
		return probeHTTP("http://127.0.0.1/", 3*time.Second)
	}

	for _, name := range names {
		url := "http://" + name + "/"
		if !probeHTTP(url, 5*time.Second) {
			return false
		}
	}
	return true
}

// extractServerNames scans file contents for server_name directives.
func extractServerNames(files []ReleaseFile) []string {
	seen := map[string]bool{}
	var names []string
	for _, f := range files {
		if f.Action == "delete" {
			continue
		}
		for _, match := range serverNameRE.FindAllStringSubmatch(f.Content, -1) {
			for _, name := range strings.Fields(match[1]) {
				if name == "_" || name == "localhost" || seen[name] {
					continue
				}
				seen[name] = true
				names = append(names, name)
			}
		}
	}
	return names
}

// probeHTTP sends a GET request and returns true when the status code is below 500.
func probeHTTP(url string, timeout time.Duration) bool {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}

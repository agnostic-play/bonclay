package api_mock_services

import (
	"regexp"
	"time"
)

const (
	requestTimeout = 30 * time.Second
)

var (
	// hopByHopHeaders are headers that should not be forwarded in proxy requests
	hopByHopHeaders = map[string]bool{
		"Connection":          true,
		"Keep-Alive":          true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true,
		"Trailers":            true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
	}

	// variableRegex matches template variables in format {{varName}}
	variableRegex = regexp.MustCompile(`\{\{(.+?)\}\}`)
)

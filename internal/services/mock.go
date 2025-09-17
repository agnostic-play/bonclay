package services

import (
	"context"
	`encoding/json`
	"fmt"
	"io"
	`log`
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	repository2 `github.com/agnostic-play/ditoo/internal/repository`
)

// MockEntityRes represents a mock API response entity
type MockEntityRes struct {
	EndpointID          string      `json:"endpoint_id" gorm:"column:endpoint_id"`
	Desc                string      `json:"desc" gorm:"column:desc"`
	Header              string      `json:"header" gorm:"column:header"`
	Body                string      `json:"body" gorm:"column:body"`
	StatusHeader        int         `json:"status_header" gorm:"column:status_header"`
	Delay               *int        `json:"delay" gorm:"column:delay"`
	ProxyIsEnabled      bool        `json:"proxy_is_enabled" gorm:"column:proxy_is_enabled"`
	ProxyResponseHeader http.Header `json:"proxy_response_header" gorm:"-"`
	ProxyResponseBody   []byte      `json:"proxy_response_body" gorm:"-"`
}

// MockServiceInterface defines the interface for mock API operations
type MockServiceInterface interface {
	MockAPI(ctx echo.Context, collectionSlug, method, path string) (MockEntityRes, error)
}

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

// MockAPI handles mock API requests with optional proxy forwarding
func (cont serviceContainer) MockAPI(echoCtx echo.Context, collectionSlug, method, apiPath string) (MockEntityRes, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/%s", apiPath)

	// Get collection by slug
	collection, err := cont.repoContainer.GetCollectionBySlug(ctx, collectionSlug)
	if err != nil {
		return MockEntityRes{}, fmt.Errorf("collection not found: %w", err)
	}

	// Get endpoint mock configuration
	endpoint, err := cont.repoContainer.GetEndpointMock(ctx, collection.ID.String(), method, path)
	if err != nil {
		return MockEntityRes{}, fmt.Errorf("failed to get endpoint mock: %w", err)
	}

	// Validate active scenario

	var mockScenario MockEntityRes

	// Handle proxy forwarding if enabled
	if collection.IsProxyEnable && collection.ForwardProxyURL != "" {
		proxyResponse, err := cont.forwardToProxy(echoCtx, collection.ForwardProxyURL, path)
		if err != nil {
			return MockEntityRes{}, fmt.Errorf("proxy forwarding failed: %w", err)
		}

		if endpoint.ActiveScenario == "" {
			return proxyResponse, nil
		}
	}

	// Get scenario configuration
	scenario, err := cont.repoContainer.GetScenario(ctx, endpoint.ActiveScenario)
	if err != nil {
		return MockEntityRes{}, fmt.Errorf("failed to get scenario: %w", err)
	}

	// Convert scenario to mock response structure
	if err := convertStruct(scenario, &mockScenario); err != nil {
		return MockEntityRes{}, fmt.Errorf("failed to convert scenario: %w", err)
	}

	// Apply custom variables
	customVars, err := cont.repoContainer.GetListCustomVariableByCollectionId(ctx, collection.ID.String())
	if err != nil {
		return MockEntityRes{}, fmt.Errorf("failed to get custom variables: %w", err)
	}

	mockScenario.applyEnvironmentVariables(customVarsToMap(customVars))

	return mockScenario, nil
}

// forwardToProxy forwards the request to a proxy server
func (cont serviceContainer) forwardToProxy(echoCtx echo.Context, proxyURL, path string) (MockEntityRes, error) {
	targetURL := fmt.Sprintf("%s%s", proxyURL, path)

	client := &http.Client{
		Timeout: requestTimeout,
	}

	// Create proxy request
	req, err := http.NewRequest(echoCtx.Request().Method, targetURL, echoCtx.Request().Body)
	if err != nil {
		log.Printf("[ERROR] failed to create proxy request: %v", err)
		return MockEntityRes{}, fmt.Errorf("failed to create proxy request: %w", err)
	}

	// Copy original request headers
	for key, values := range echoCtx.Request().Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	marshal, err := json.Marshal(echoCtx.Request().Body)
	if err != nil {
		return MockEntityRes{}, err
	}

	log.Printf("[REQUEST] Forwarding request: method=%s url=%s headers=%v body=%v",
		req.Method, req.URL.String(), req.Header, marshal)

	// Execute proxy request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] proxy request failed: %v", err)
		return MockEntityRes{}, fmt.Errorf("proxy request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] failed to read proxy response body: %v", err)
		return MockEntityRes{}, fmt.Errorf("failed to read proxy response body: %w", err)
	}

	log.Printf("[RESPONSE] status=%d headers=%v body=%s",
		resp.StatusCode, resp.Header, string(bodyBytes))

	// Prepare proxy response
	proxyHeaders := make(http.Header)
	copyHeaders(resp.Header, proxyHeaders)

	return MockEntityRes{
		ProxyIsEnabled:      true,
		StatusHeader:        resp.StatusCode,
		ProxyResponseHeader: proxyHeaders,
		ProxyResponseBody:   bodyBytes,
	}, nil
}

// applyEnvironmentVariables replaces template variables in the response body
func (res *MockEntityRes) applyEnvironmentVariables(env map[string]string) {
	res.Body = variableRegex.ReplaceAllStringFunc(res.Body, func(match string) string {
		// Extract variable name without {{ }}
		submatch := variableRegex.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		varName := strings.TrimSpace(submatch[1])
		if value, exists := env[varName]; exists {
			return value
		}

		// Return original placeholder if variable not found
		return match
	})
}

// customVarsToMap converts a slice of custom variables to a map
func customVarsToMap(customVars []repository2.CustomVariableEntity) map[string]string {
	result := make(map[string]string, len(customVars))
	for _, variable := range customVars {
		result[variable.Key] = variable.Value
	}
	return result
}

// convertStruct copies fields from source struct to destination struct using reflection
func convertStruct(src, dst interface{}) error {
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.ValueOf(dst)

	// Validate inputs
	if srcValue.Kind() != reflect.Struct {
		return fmt.Errorf("source must be a struct")
	}
	if dstValue.Kind() != reflect.Ptr || dstValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to struct")
	}

	dstElem := dstValue.Elem()
	srcType := srcValue.Type()

	// Copy matching fields
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		fieldName := srcType.Field(i).Name

		dstField := dstElem.FieldByName(fieldName)
		if dstField.IsValid() && dstField.CanSet() && srcField.Type() == dstField.Type() {
			dstField.Set(srcField)
		}
	}

	return nil
}

// copyHeaders copies HTTP headers from source to destination, excluding hop-by-hop headers
func copyHeaders(src, dst http.Header) {
	for key, values := range src {
		if !hopByHopHeaders[key] {
			for _, value := range values {
				dst.Add(key, value)
			}
		}
	}
}

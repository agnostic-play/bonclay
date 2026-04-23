package api_mock_services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	repository2 `berlin.allobank.com/tools/bonclay/internal/adapters/repositories`
	"github.com/labstack/echo/v4"

	"berlin.allobank.com/tools/bonclay/internal/common/pagination"
	"berlin.allobank.com/tools/bonclay/internal/entities"
	crud_services "berlin.allobank.com/tools/bonclay/internal/services/crud-services"
)

// MockService defines the interface for mock API operations
type MockService interface {
	MockAPI(ctx echo.Context, collectionSlug, method, path string) (MockEntityRes, error)
}

type mockService struct {
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity]
	scenarioService   crud_services.BaseCRUDService[entities.ScenarioEntity]
	repoContainer     repository2.RepoContainer
}

func NewMockService(
	collectionService crud_services.BaseCRUDService[entities.CollectionEntity],
	scenarioService crud_services.BaseCRUDService[entities.ScenarioEntity],
	repoContainer repository2.RepoContainer,
) MockService {
	return &mockService{
		collectionService: collectionService,
		scenarioService:   scenarioService,
		repoContainer:     repoContainer,
	}
}

// MockAPI handles mock API requests with optional proxy forwarding
func (s *mockService) MockAPI(echoCtx echo.Context, collectionSlug, method, apiPath string) (MockEntityRes, error) {
	ctx := context.Background()
	path := fmt.Sprintf("/%s", apiPath)

	// Get collection by slug using CRUD service
	collQuery := &pagination.ListQuery{ShowAll: true, Filters: map[string]interface{}{"slug = ?": collectionSlug}}
	collResult, err := s.collectionService.GetList(ctx, collQuery)
	if err != nil || len(collResult.List) == 0 {
		return MockEntityRes{}, fmt.Errorf("collection not found")
	}
	collection := collResult.List[0]

	// Get endpoint mock configuration using repository
	listEndpoint, err := s.repoContainer.GetEndpointRepo().GetListEndpointByCollectionID(ctx, collection.ID.String(), method)
	if err != nil {
		return MockEntityRes{}, fmt.Errorf("failed to get endpoint mock: %w", err)
	}

	if listEndpoint == nil || len(listEndpoint) < 1 {
		return MockEntityRes{}, fmt.Errorf("endpoint mock not found")
	}

	var matchedEndpoint *entities.EndpointEntity
	for _, data := range listEndpoint {
		re := regexp.MustCompile("^" + data.Path + "$")
		if re.MatchString(path) {
			matchedEndpoint = &data
			break
		}
	}

	if matchedEndpoint == nil {
		return MockEntityRes{}, fmt.Errorf("endpoint mock not found")
	}

	endpoint := matchedEndpoint
	var mockScenario MockEntityRes

	// Handle proxy forwarding if enabled
	if collection.IsProxyEnable != nil && *collection.IsProxyEnable && collection.ForwardProxyURL != "" && endpoint.ActiveScenario == "" {
		proxyResponse, err := s.forwardToProxy(echoCtx, *endpoint, collection.ForwardProxyURL, path)
		if err != nil {
			return MockEntityRes{}, fmt.Errorf("proxy forwarding failed: %w", err)
		}

		mockScenario = proxyResponse
	}

	if endpoint.ActiveScenario != "" {
		// Get scenario using CRUD service
		scenario, err := s.scenarioService.Get(ctx, endpoint.ActiveScenario)
		if err != nil {
			return MockEntityRes{}, fmt.Errorf("failed to get scenario: %w", err)
		}

		// Convert scenario to mock response structure
		if err := convertStruct(scenario, &mockScenario); err != nil {
			return MockEntityRes{}, fmt.Errorf("failed to convert scenario: %w", err)
		}
	}

	// Apply custom variables using repository
	customVars, err := s.repoContainer.GetCustomVariableRepo().GetListCustomVariableByCollectionId(ctx, collection.ID.String())
	if err != nil {
		return MockEntityRes{}, fmt.Errorf("failed to get custom variables: %w", err)
	}

	mockScenario.applyEnvironmentVariables(customVarsToMap(customVars))

	return mockScenario, nil
}

// forwardToProxy forwards the request to a proxy server
func (s *mockService) forwardToProxy(echoCtx echo.Context, endpoint entities.EndpointEntity, proxyURL, path string) (MockEntityRes, error) {
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

	if echoCtx.Request().Body != nil {
		marshal, _ := json.Marshal(echoCtx.Request().Body)
		log.Printf("[REQUEST] Forwarding request: method=%s url=%s headers=%v body=%v",
			req.Method, req.URL.String(), req.Header, marshal)
	}

	// Execute proxy request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] proxy request failed: %v", err)
		return MockEntityRes{}, fmt.Errorf("proxy request failed: %w", err)
	}

	defer resp.Body.Close()
	var bodyBytes []byte
	if resp.Body != nil {
		// Read response body
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR] failed to read proxy response body: %v", err)
			return MockEntityRes{}, fmt.Errorf("failed to read proxy response body: %w", err)
		}

		log.Printf("[RESPONSE] status=%d headers=%v body=%s",
			resp.StatusCode, resp.Header, string(bodyBytes))
	}

	// Prepare proxy response
	proxyHeaders := make(http.Header)
	copyHeaders(resp.Header, proxyHeaders)

	return MockEntityRes{
		EndpointID:          endpoint.ID.String(),
		ProxyIsEnabled:      true,
		StatusHeader:        resp.StatusCode,
		Delay:               endpoint.Delay,
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

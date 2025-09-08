package services

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/agnostic-play/ditoo/internal/repository"
)

type MockEntityRes struct {
	EndpointId   string `json:"endpoint_id" gorm:"column:endpoint_id"`
	Desc         string `json:"desc" gorm:"column:desc"`
	Header       string `json:"header" gorm:"column:header"`
	Body         string `json:"body" gorm:"column:body"`
	StatusHeader int    `json:"status_header" gorm:"column:status_header"`
	Delay        *int   `json:"delay" gorm:"column:delay"`
}

type MockServiceInterface interface {
	MockApi(ctx context.Context, collectionSlug, method string, path string) (MockEntityRes, error)
}

func (cont serviceContainer) MockApi(ctx context.Context, collectionSlug, method string, path string) (MockEntityRes, error) {
	path = fmt.Sprintf("/%s", path)
	collection, err := cont.repoContainer.GetCollectionBySlug(ctx, collectionSlug)
	if err != nil {
		err = fmt.Errorf("collection not found")
		return MockEntityRes{}, err
	}

	endpoint, err := cont.repoContainer.GetEndpointMock(ctx, collection.ID.String(), method, path)
	if err != nil {
		return MockEntityRes{}, err
	}

	if strings.TrimSpace(endpoint.ActiveScenario) == "" {
		err = fmt.Errorf("scenario not found")
		return MockEntityRes{}, err
	}

	var mockScenario MockEntityRes
	scenario, err := cont.repoContainer.GetScenario(ctx, endpoint.ActiveScenario)
	if err != nil {
		return MockEntityRes{}, err
	}
	ConvertStruct(scenario, &mockScenario)

	customVar, err := cont.repoContainer.GetListCustomVariableByCollectionId(ctx, collection.ID.String())
	if err != nil {
		return MockEntityRes{}, err
	}
	mockScenario.ApplyEnv(envSliceToMap(customVar))

	return mockScenario, nil
}

func (scenario *MockEntityRes) ApplyEnv(env map[string]string) {
	// Match {{varName}}
	re := regexp.MustCompile(`\{\{(.+?)\}\}`)

	// Replace each match with env value
	result := re.ReplaceAllStringFunc(scenario.Body, func(m string) string {
		// Extract key name without {{ }}
		key := re.FindStringSubmatch(m)[1]
		if val, ok := env[key]; ok {
			return val
		}
		// If not found, keep original placeholder
		return m
	})

	scenario.Body = result

	return
}

func envSliceToMap(envs []repository.CustomVariableEntity) map[string]string {
	result := make(map[string]string)
	for _, e := range envs {
		result[e.Key] = e.Value
	}
	return result
}

func ConvertStruct(src, dst interface{}) {
	sv := reflect.ValueOf(src)
	dv := reflect.ValueOf(dst).Elem()

	for i := 0; i < sv.NumField(); i++ {
		field := sv.Type().Field(i).Name
		if f := dv.FieldByName(field); f.IsValid() && f.CanSet() {
			f.Set(sv.Field(i))
		}
	}
}

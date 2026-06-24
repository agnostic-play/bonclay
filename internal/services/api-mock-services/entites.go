package api_mock_services

import "net/http"

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

type CustomVariableReq struct {
	Method       string `json:"method" validate:"required,alpha"`
	CollectionID string `json:"collection_id" validate:"required,uuid"`
	Key          string `json:"key" validate:"required"`
	Value        string `json:"value" validate:"required"`
}

// CollectionVariableReq is the clean v2 request for creating/updating a variable.
type CollectionVariableReq struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type SetActiveScenarioReq struct {
	EndpointID string `json:"endpoint_id" validate:"required,uuid"`
	ScenarioID string `json:"scenario_id" validate:"required"`
}

// EndpointScriptReq is the request body for updating an endpoint's mock script.
// Script may be empty to clear an existing script.
type EndpointScriptReq struct {
	Script string `json:"script"`
}

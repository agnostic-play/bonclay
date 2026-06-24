package api_mock_services

import (
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
)

// ScriptRequest is the read-only view of the incoming mock request that is
// exposed to a script as the global `request` object.
type ScriptRequest struct {
	Method  string
	Path    string
	Body    string
	Headers map[string]string
}

// ScriptResult is the outcome of running an endpoint script: the (possibly
// mutated) variables plus any response directives the script set on the global
// `response` object.
type ScriptResult struct {
	Env map[string]string

	// ScenarioID, when set via `response.scenario`, selects which saved scenario
	// to serve for this request, overriding the endpoint's active scenario.
	ScenarioID string
	// Body / Status, when set via `response.body` / `response.status`, override
	// the resolved scenario's body and status code. nil means "not set".
	Body   *string
	Status *int
}

// RunScript executes a user-provided JavaScript snippet for a mock request.
//
// Three globals are available to the script:
//
//   - env      — a mutable object of the collection's custom variables. Reading
//     and assigning keys (e.g. env["token"] = "abc") updates the stored values.
//   - request  — the incoming request: { method, path, body, json, headers }.
//     `request.json` is the parsed body when it is valid JSON, otherwise null.
//   - response — a writable object the script can use to drive the reply:
//     `response.scenario` (pick a saved scenario by id), `response.body`
//     (override the body), `response.status` (override the status code).
//
// If the script fails to execute, the original env is returned with the error so
// callers can decide whether to proceed with the unmodified variables.
func RunScript(env map[string]string, script string, req ScriptRequest) (ScriptResult, error) {
	result := ScriptResult{Env: env}

	vm := goja.New()

	// Expose env as a mutable object.
	jsEnv := vm.NewObject()
	for key, value := range env {
		_ = jsEnv.Set(key, value)
	}
	if err := vm.Set("env", jsEnv); err != nil {
		return result, fmt.Errorf("failed to expose env to script: %w", err)
	}

	// Expose the incoming request.
	jsReq := vm.NewObject()
	_ = jsReq.Set("method", req.Method)
	_ = jsReq.Set("path", req.Path)
	_ = jsReq.Set("body", req.Body)

	jsHeaders := vm.NewObject()
	for key, value := range req.Headers {
		_ = jsHeaders.Set(key, value)
	}
	_ = jsReq.Set("headers", jsHeaders)

	// Pre-parse JSON bodies for convenience as request.json (null otherwise).
	if req.Body != "" {
		var parsed interface{}
		if err := json.Unmarshal([]byte(req.Body), &parsed); err == nil {
			_ = jsReq.Set("json", vm.ToValue(parsed))
		} else {
			_ = jsReq.Set("json", goja.Null())
		}
	} else {
		_ = jsReq.Set("json", goja.Null())
	}
	_ = vm.Set("request", jsReq)

	// Expose a writable response object for the script to drive output.
	jsResp := vm.NewObject()
	if err := vm.Set("response", jsResp); err != nil {
		return result, fmt.Errorf("failed to expose response to script: %w", err)
	}

	if _, err := vm.RunString(script); err != nil {
		return result, fmt.Errorf("script execution failed: %w", err)
	}

	// Read the (possibly mutated) env back out.
	updated := make(map[string]string, len(jsEnv.Keys()))
	for _, key := range jsEnv.Keys() {
		value := jsEnv.Get(key)
		if value == nil {
			continue
		}
		if str, ok := value.Export().(string); ok {
			updated[key] = str
		} else {
			updated[key] = fmt.Sprintf("%v", value.Export())
		}
	}
	result.Env = updated

	// Read back any response directives the script set.
	if v := jsResp.Get("scenario"); isSet(v) {
		result.ScenarioID = v.String()
	}
	if v := jsResp.Get("body"); isSet(v) {
		body := v.String()
		result.Body = &body
	}
	if v := jsResp.Get("status"); isSet(v) {
		status := int(v.ToInteger())
		result.Status = &status
	}

	return result, nil
}

// isSet reports whether a script-provided value is present (not nil/undefined/null).
func isSet(v goja.Value) bool {
	return v != nil && !goja.IsUndefined(v) && !goja.IsNull(v)
}

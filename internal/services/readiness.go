package services

type ReadinessResp struct {
	Dependencies map[string]string `json:"dependencies"`
}

func (cont serviceContainer) HelloWorld() ReadinessResp {
	resp := ReadinessResp{}
	resp.Dependencies["database"] = "ready"
	if cont.repoContainer.Ping() != nil {
		resp.Dependencies["database"] = "not ready"
	}

	return resp
}

func (cont serviceContainer) Readiness() ReadinessResp {
	var resp ReadinessResp

	resp.Dependencies = make(map[string]string)

	resp.Dependencies["database"] = "ready"
	if cont.repoContainer.Ping() != nil {
		resp.Dependencies["database"] = "not ready"
	}
	return resp
}

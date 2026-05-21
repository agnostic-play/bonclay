package kong_services

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"berlin.allobank.com/tools/bonclay/internal/config"
)

type KongService interface {
	ListServices(ctx context.Context, query ListServicesQuery) (*ListServicesResponse, error)
	UpdateService(ctx context.Context, id string, patch UpdateServicePayload) (*Service, error)
}

type kongService struct {
	cfg    config.Kong
	client *http.Client
}

func NewKongService(cfg config.Kong) KongService {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.Insecure,
		},
	}

	return &kongService{
		cfg: cfg,
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

type ListServicesQuery struct {
	Size   int
	Offset string
	Tag    string
}

type Service struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Protocol          string   `json:"protocol"`
	Host              string   `json:"host"`
	Port              int      `json:"port"`
	Path              *string  `json:"path"`
	Tags              []string `json:"tags"`
	Enabled           bool     `json:"enabled"`
	Retries           int      `json:"retries"`
	ConnectTimeout    int      `json:"connect_timeout"`
	ReadTimeout       int      `json:"read_timeout"`
	WriteTimeout      int      `json:"write_timeout"`
	ClientCertificate *string  `json:"client_certificate"`
	CACertificates    []string `json:"ca_certificates"`
	TLSVerify         *bool    `json:"tls_verify"`
	TLSVerifyDepth    *int     `json:"tls_verify_depth"`
	CreatedAt         int64    `json:"created_at"`
	UpdatedAt         int64    `json:"updated_at"`
}

type ListServicesResponse struct {
	Data []Service `json:"data"`
	Next *string   `json:"next"`
}

type UpdateServicePayload struct {
	Enabled *bool `json:"enabled,omitempty"`
}

func (s *kongService) ListServices(ctx context.Context, query ListServicesQuery) (*ListServicesResponse, error) {
	endpoint, err := s.buildURL("/services")
	if err != nil {
		return nil, err
	}

	q := endpoint.Query()
	if query.Size > 0 {
		q.Set("size", fmt.Sprintf("%d", query.Size))
	}
	if strings.TrimSpace(query.Offset) != "" {
		q.Set("offset", query.Offset)
	}
	if strings.TrimSpace(query.Tag) != "" {
		q.Set("tags", query.Tag)
	}
	endpoint.RawQuery = q.Encode()

	body, err := s.do(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	var out ListServicesResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("failed parsing kong response: %w", err)
	}

	return &out, nil
}

func (s *kongService) UpdateService(ctx context.Context, id string, patch UpdateServicePayload) (*Service, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("service id is required")
	}

	endpoint, err := s.buildURL("/services/" + url.PathEscape(id))
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(patch)
	if err != nil {
		return nil, fmt.Errorf("failed encoding patch payload: %w", err)
	}

	body, err := s.do(ctx, http.MethodPatch, endpoint.String(), payload)
	if err != nil {
		return nil, err
	}

	var out Service
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("failed parsing kong response: %w", err)
	}

	return &out, nil
}

func (s *kongService) buildURL(p string) (*url.URL, error) {
	if strings.TrimSpace(s.cfg.AdminURL) == "" {
		return nil, fmt.Errorf("kong adminURL is not configured")
	}
	base := strings.TrimRight(s.cfg.AdminURL, "/")
	u, err := url.Parse(base + p)
	if err != nil {
		return nil, fmt.Errorf("invalid kong adminURL: %w", err)
	}
	return u, nil
}

func (s *kongService) do(ctx context.Context, method, fullURL string, body []byte) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(s.cfg.AdminToken) != "" {
		req.Header.Set("Kong-Admin-Token", s.cfg.AdminToken)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed calling kong admin: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading kong response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("kong admin returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

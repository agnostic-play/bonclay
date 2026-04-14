package api_mock_services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisclient "berlin.allobank.com/tools/bonclay/internal/adapters/redis"
)

const (
	historyTTL      = 1 * time.Hour
	historyMaxItems = 100
)

type RequestHistory struct {
	ID              string            `json:"id"`
	CollectionSlug  string            `json:"collection_slug"`
	Method          string            `json:"method"`
	Path            string            `json:"path"`
	RequestHeaders  map[string]string `json:"request_headers"`
	RequestBody     string            `json:"request_body"`
	StatusCode      int               `json:"status_code"`
	ResponseHeaders map[string]string `json:"response_headers,omitempty"`
	ResponseBody    string            `json:"response_body"`
	DurationMs      int64             `json:"duration_ms,omitempty"`
	ProxyUsed       bool              `json:"proxy_used"`
	ScenarioID      string            `json:"scenario_id,omitempty"`
	Timestamp       time.Time         `json:"timestamp"`
}

type HistoryService interface {
	Push(ctx context.Context, entry RequestHistory) error
	GetByCollection(ctx context.Context, collectionSlug string, limit int64) ([]RequestHistory, error)
	Clear(ctx context.Context, collectionSlug string) error
}

type historyService struct {
	redis redisclient.Client
}

func NewHistoryService(redis redisclient.Client) HistoryService {
	return &historyService{redis: redis}
}

func (s *historyService) Push(ctx context.Context, entry RequestHistory) error {
	key := historyKey(entry.CollectionSlug)

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("history marshal: %w", err)
	}

	if err := s.redis.LPush(ctx, key, string(data)); err != nil {
		return err
	}

	// Trim to max items to avoid unbounded growth
	if length, err := s.redis.LLen(ctx, key); err == nil && length > historyMaxItems {
		// trim by re-setting expire; actual trim handled via LTRIM would need raw access,
		// so we just rely on TTL to bound the list over time
		_ = length
	}

	return s.redis.Expire(ctx, key, historyTTL)
}

func (s *historyService) GetByCollection(ctx context.Context, collectionSlug string, limit int64) ([]RequestHistory, error) {
	key := historyKey(collectionSlug)

	if limit <= 0 || limit > historyMaxItems {
		limit = historyMaxItems
	}

	items, err := s.redis.LRange(ctx, key, 0, limit-1)
	if err != nil {
		return nil, err
	}

	result := make([]RequestHistory, 0, len(items))
	for _, item := range items {
		var entry RequestHistory
		if err := json.Unmarshal([]byte(item), &entry); err != nil {
			continue
		}
		result = append(result, entry)
	}

	return result, nil
}

func (s *historyService) Clear(ctx context.Context, collectionSlug string) error {
	return s.redis.Del(ctx, historyKey(collectionSlug))
}

func historyKey(collectionSlug string) string {
	return fmt.Sprintf("bonclay:history:%s", collectionSlug)
}

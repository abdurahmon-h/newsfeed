package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rssparser/internal/models"
)

type SearchResponse struct {
	Hits struct {
		Hits []struct {
			Source models.NewsItems `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func SearchNews(keyword string) error {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": keyword,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return fmt.Errorf("ошибка кодирования запроса: %w", err)
	}

	resp, err := http.Post("http://localhost:9200/news/_search", "application/json", &buf)
	if err != nil {
		return fmt.Errorf("ошибка выполнения поиска: %w", err)
	}
	defer resp.Body.Close()

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return fmt.Errorf("ошибка разбора ответа поиска: %w", err)
	}

	for _, hit := range searchResp.Hits.Hits {
		fmt.Println("🔍 Найдено:", hit.Source.Title)
	}

	return nil
}

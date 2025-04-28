package storage

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"rssparser/internal/models"
)

type SearchResponse struct { // вложенная структура ес
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
	if err := json.NewEncoder(&buf).Encode(query); err != nil { // кодирование запроса в json
		log.Printf("ошибка кодирования запроса: %v", err)
		return err
	}

	resp, err := http.Post("http://localhost:9200/news/_search", "application/json", &buf) // отправляем запрос в ес
	if err != nil {
		log.Printf("ошибка выполнения поиска: %w", err)
		return err
	}
	defer resp.Body.Close()

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		log.Printf("ошибка разбора ответа поиска: %w", err)
		return err
	}

	for _, hit := range searchResp.Hits.Hits {
		log.Println("Найдено:", hit.Source.Title)
	}

	return nil
}

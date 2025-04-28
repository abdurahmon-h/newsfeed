package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"

	"rssparser/internal/models"
)

type ElasticStorage struct {
	Client *elasticsearch.Client
	Index  string
}

func NewElasticStorage(indexName string) (*ElasticStorage, error) {
	conf := elasticsearch.Config{ // создаем конфиг для ес
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	esClient, err := elasticsearch.NewClient(conf) // создаем клиеента для ес
	if err != nil {
		return nil, err
	}

	return &ElasticStorage{
		Client: esClient,
		Index:  indexName,
	}, nil
}

func (e *ElasticStorage) SaveNewsItem(item models.NewsItems) error {
	data, err := json.Marshal(item) // конвертация айтемов в json
	if err != nil {
		return err
	}

	id := uuid.New().String() // генерация ид для каждого айтема
	res, err := e.Client.Index(
		e.Index,
		bytes.NewReader(data),
		e.Client.Index.WithDocumentID(id),
		e.Client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("error indexing document: %s", res.String())
		return err
	}

	return nil
}

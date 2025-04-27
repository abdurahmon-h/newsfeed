package parser

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"rssparser/internal/models"
)

type RSSParser struct {
	FeedURL string
	Source  string
}

type rss struct {
	Channel channel `xml:"channel"`
}

type channel struct {
	Items []item `xml:"item"`
}

type item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

func (p *RSSParser) GetNews() ([]models.NewsItems, error) { // функция для парсинга рсс фидов
	resp, err := http.Get(p.FeedURL)
	if err != nil {
		log.Printf("failed to get RSS: %w", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read RSS: %w", err)
		return nil, err
	}

	var parse rss
	if err := xml.Unmarshal(body, &parse); err != nil {
		log.Printf("failed to parse RSS: %w", err)
		return nil, err
	}

	var news []models.NewsItems
	for _, i := range parse.Channel.Items {
		pubTime, err := parsePubDate(i.PubDate)
		if err != nil {
			pubTime = time.Now()
		}

		news = append(news, models.NewsItems{ // цбираем пробелы в текстах
			Title:       strings.TrimSpace(i.Title),
			Description: strings.TrimSpace(i.Description),
			Link:        strings.TrimSpace(i.Link),
			Source:      p.Source,
			PublishedAt: pubTime,
		})
	}

	return news, nil
}

func parsePubDate(dateStr string) (time.Time, error) {
	timeModel := time.RFC1123 // принимаем дату и переводим его в формат дат рсс
	return time.Parse(timeModel, dateStr)
}

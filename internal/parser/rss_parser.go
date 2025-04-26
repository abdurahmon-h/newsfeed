package parser

import (
	"encoding/xml"
	"fmt"
	"io"
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

func (p *RSSParser) FetchNews() ([]models.NewsItems, error) {
	resp, err := http.Get(p.FeedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RSS body: %w", err)
	}

	var parsed rss
	if err := xml.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	var news []models.NewsItems
	for _, i := range parsed.Channel.Items {
		pubTime, err := parsePubDate(i.PubDate)
		if err != nil {
			pubTime = time.Now()
		}

		news = append(news, models.NewsItems{
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
	layout := time.RFC1123
	return time.Parse(layout, dateStr)
}

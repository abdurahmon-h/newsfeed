package main

import (
	"fmt"
	"log"
	"rssparser/internal/parser"
	"rssparser/internal/storage"
)

func main() {
	habr := parser.RSSParser{
		FeedURL: "https://habr.com/ru/rss/all/all/?fl=ru",
		Source:  "Habr",
	}

	news, err := habr.GetNews()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Получено %d новостей", len(news))

	elastic, err := storage.NewElasticStorage("news")
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range news {
		if err := elastic.SaveNewsItem(item); err != nil {
			log.Println("Ошибка при сохранении:", err)
		} else {
			log.Println("Сохранено:", item.Title)
		}
	}

	query := "такое" // поиск запросов по ключевому слову
	if err := storage.SearchNews(query); err != nil {
		fmt.Println("Ошибка при поиске:", err)
	}
}

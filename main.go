package main

import (
	"fmt"
	"rssparser/internal/parser"
	"rssparser/internal/storage"
)

func main() {
	habr := parser.RSSParser{
		FeedURL: "https://habr.com/ru/rss/all/all/?fl=ru",
		Source:  "Habr",
	}

	news, err := habr.FetchNews()
	if err != nil {
		panic(err)
	}

	fmt.Printf("✅ Получено %d новостей\n", len(news))

	elastic, err := storage.NewElasticStorage("news")
	if err != nil {
		panic(err)
	}

	for _, item := range news {
		if err := elastic.SaveNewsItem(item); err != nil {
			fmt.Println("❌ Ошибка при сохранении:", err)
		} else {
			fmt.Println("💾 Сохранено:", item.Title)
		}
	}

	query := "Как" // поиск запросов по ключевому слову
	if err := storage.SearchNews(query); err != nil {
		fmt.Println("❌ Ошибка при поиске:", err)
	}
}

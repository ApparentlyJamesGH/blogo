package articles

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/pluja/blogo/internal/models"
)

var ArticleList []models.Article

func WatchArticles() {
	articlesPath := viper.GetString("articlesPath")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					updateArticleList(articlesPath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(articlesPath)
	if err != nil {
		log.Fatal(err)
	}

	updateArticleList(articlesPath)

	select {}
}

func updateArticleList(articlesPath string) {
	files, err := os.ReadDir(articlesPath)
	if err != nil {
		log.Println("Error reading directory:", err)
		return
	}

	var newList []models.Article
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			filename := strings.TrimSuffix(file.Name(), ".md")
			article, err := GetFromFile(filename, false)
			if err != nil {
				log.Printf("Error reading article %s: %v\n", filename, err)
				continue
			}
			if !article.Draft {
				newList = append(newList, article)
			}
		}
	}

	sort.Slice(newList, func(i, j int) bool {
		return newList[i].Date.After(newList[j].Date)
	})

	ArticleList = newList
}

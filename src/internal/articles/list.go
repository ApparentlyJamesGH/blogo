package articles

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/pluja/blogo/internal/cache"
	"github.com/pluja/blogo/internal/models"
	"github.com/pluja/blogo/internal/nostr"
)

var (
	ArticleMap        sync.Map
	ArticleList       []models.Article
	TagMap            = make(map[string][]string)
	mutex             sync.Mutex
	nostrPublishTimer *time.Timer
	nostrPublishMutex sync.Mutex
)

const nostrPublishDelay = 5 * time.Minute

func WatchArticles() {
	articles_path := viper.GetString("articles_path")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Err(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				filename := filepath.Base(event.Name)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Debug().Msgf("Article updated.")
					updateArticle(filename)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Debug().Msgf("Article removed.")
					removeArticle(filename)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Debug().Msgf("Article renamed.")
					removeArticle(filename)
				}
				UpdateFeed()
				cache.Cache.Del(strings.TrimSuffix(filename, ".md"))
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error().Err(err).Msg("watcher error")
			}
		}
	}()

	err = watcher.Add(articles_path)
	if err != nil {
		log.Error().Err(err).Msg("watcher error")
	}

	updateArticleMap(articles_path)
	UpdateFeed()

	select {}
}

func updateArticleMap(articles_path string) {
	files, err := os.ReadDir(articles_path)
	if err != nil {
		log.Error().Err(err).Msg("Error reading directory:")
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	ArticleList = nil
	TagMap = make(map[string][]string)

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			filename := strings.TrimSuffix(file.Name(), ".md")
			if filename == "about" {
				continue
			}
			article, err := GetFromFile(filename, false)
			if err != nil {
				log.Printf("Error reading article %s: %v\n", filename, err)
				continue
			}
			if !article.Draft {
				ArticleList = append(ArticleList, article)
				updateTagMap(article)
			}
		}
	}

	sort.Slice(ArticleList, func(i, j int) bool {
		return ArticleList[i].Date.After(ArticleList[j].Date)
	})

	ArticleMap = sync.Map{}
	for _, article := range ArticleList {
		ArticleMap.Store(article.Slug, article)
	}

	if viper.GetBool("nostr.publish") {
		scheduleNostrPublish()
	}
}

func updateArticle(filename string) {
	article, err := GetFromFile(strings.TrimSuffix(filename, ".md"), false)
	if err != nil {
		log.Printf("Error updating article %s: %v\n", filename, err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if !article.Draft {
		ArticleMap.Store(article.Slug, article)
		updateTagMap(article)
	} else {
		ArticleMap.Delete(article.Slug)
		updateTagMapOnRemoval(article)
	}

	updateArticleList()
	if viper.GetBool("nostr.publish") {
		scheduleNostrPublish()
	}
}

func removeArticle(filename string) {
	mutex.Lock()
	defer mutex.Unlock()

	articleSlug := strings.TrimSuffix(filename, ".md")
	if value, ok := ArticleMap.Load(articleSlug); ok {
		article := value.(models.Article)
		updateTagMapOnRemoval(article)
	}
	ArticleMap.Delete(articleSlug)
	updateArticleMap(viper.GetString("articles_path"))

	if viper.GetBool("nostr.publish") {
		scheduleNostrPublish()
	}
}

func updateTagMap(article models.Article) {
	for _, tag := range article.Tags {
		TagMap[tag] = append(TagMap[tag], article.Slug)
	}
}

func updateTagMapOnRemoval(article models.Article) {
	for _, tag := range article.Tags {
		slugList := TagMap[tag]
		for i, slug := range slugList {
			if slug == article.Slug {
				TagMap[tag] = append(slugList[:i], slugList[i+1:]...)
				break
			}
		}
	}
}

func updateArticleList() {
	ArticleList = nil

	ArticleMap.Range(func(_, value interface{}) bool {
		article := value.(models.Article)
		ArticleList = append(ArticleList, article)
		return true
	})

	sort.Slice(ArticleList, func(i, j int) bool {
		return ArticleList[i].Date.After(ArticleList[j].Date)
	})
}

func scheduleNostrPublish() {
	nostrPublishMutex.Lock()
	defer nostrPublishMutex.Unlock()

	if nostrPublishTimer != nil {
		nostrPublishTimer.Stop()
	}

	nostrPublishTimer = time.AfterFunc(nostrPublishDelay, publishBlogToNostr)
}

func publishBlogToNostr() {
	nostrPublishMutex.Lock()
	defer nostrPublishMutex.Unlock()

	log.Debug().Msgf("Publishing blog to nostr...")
	for _, a := range ArticleList {
		if _, ok := nostr.NostrPublications[a.Slug]; ok {
			log.Debug().Msgf("Skip %s, already published", a.Slug)
			continue
		}

		if a.Draft || !a.Nostr {
			log.Debug().Msgf("Skip %s, draft or nostr disabled", a.Slug)
			continue
		}

		art, err := GetFromFile(a.Slug, true)
		if err != nil {
			log.Error().Err(err).Msg("error!")
			return
		}

		err = nostr.PublishArticle(art)
		if err != nil {
			log.Error().Err(err).Msg("error!")
			return
		}
	}
}

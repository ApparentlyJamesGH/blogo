// src/internal/articles/feeds.go
package articles

import (
	"fmt"
	"time"

	"github.com/gorilla/feeds"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var feed feeds.Feed

func UpdateFeed() error {
	now := time.Now()
	feed = feeds.Feed{
		Title:       viper.GetString("title"),
		Link:        &feeds.Link{Href: fmt.Sprintf("%v/rss", viper.GetString("base_url"))},
		Description: viper.GetString("description"),
		Author:      &feeds.Author{Name: viper.GetString("title")},
		Created:     now,
		Items:       getFeedItems(), // Extract feed items generation
	}
	return nil
}

func getFeedItems() []*feeds.Item {
	items := []*feeds.Item{}
	for _, article := range ArticleList {
		if !article.Draft {
			item := &feeds.Item{
				Title:       article.Title,
				Link:        &feeds.Link{Href: fmt.Sprintf("%v/p/%v", viper.GetString("base_url"), article.Slug)},
				Description: article.Summary,
				Created:     article.Date,
			}
			items = append(items, item)
		}
	}
	return items
}

func RssFeed() string {
	rss, err := feed.ToRss()
	if err != nil {
		log.Err(err).Msg("Error generating RSS feed")
		return ""
	}
	return rss
}

func AtomFeed() string {
	atom, err := feed.ToAtom()
	if err != nil {
		log.Err(err).Msg("Error generating Atom feed")
		return ""
	}
	return atom
}

func JsonFeed() string {
	json, err := feed.ToJSON()
	if err != nil {
		log.Err(err).Msg("Error generating JSON feed")
		return ""
	}
	return json
}

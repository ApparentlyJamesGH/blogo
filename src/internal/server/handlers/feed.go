// src/internal/server/handlers/feed.go
package handlers

import (
	"github.com/kataras/iris/v12"

	"github.com/pluja/blogo/internal/articles"
)

func HandleFeed(c iris.Context) {
	switch c.Path() {
	case "/rss":
		c.Header("Content-Type", "application/rss+xml")
		c.StatusCode(iris.StatusOK)
		c.WriteString(articles.RssFeed())
	case "/atom":
		c.Header("Content-Type", "application/atom+xml")
		c.StatusCode(iris.StatusOK)
		c.WriteString(articles.AtomFeed())
	case "/json":
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.StatusCode(iris.StatusOK)
		c.WriteString(articles.JsonFeed())
	default:
		c.StatusCode(iris.StatusNotFound)
		c.Text("Feed not found")
	}
}

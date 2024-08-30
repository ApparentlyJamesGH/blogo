package handlers

import (
	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog/log"

	"github.com/pluja/blogo/frontend/templates"
	"github.com/pluja/blogo/internal/articles"
	"github.com/pluja/blogo/internal/cache"
	"github.com/pluja/blogo/internal/models"
)

func HandleArticle(c iris.Context) {
	slug := c.Params().Get("slug")

	if value, found := cache.Cache.Get(slug); found {
		c.RenderComponent(templates.Article(value.(models.Article)))
		return
	}

	article, err := articles.GetFromFile(slug, true)
	if err != nil {
		log.Error().Err(err).Str("slug", slug).Msg("Error fetching article")
		c.RenderComponent(templates.ArticleNotFound(slug))
		return
	}

	cache.Cache.Set(slug, article, 1)
	c.RenderComponent(templates.Article(article))
}

func HandleRawArticle(c iris.Context) {
	slug := c.Params().Get("slug")

	if rawContent, found := cache.Cache.Get(slug + "_raw"); found {
		c.Text(rawContent.(string))
		return
	}

	article, err := articles.GetFromFile(slug, true)
	if err != nil {
		log.Error().Err(err).Str("slug", slug).Msg("Error fetching raw article")
		c.RenderComponent(templates.ArticleNotFound(slug))
		return
	}

	cache.Cache.Set(slug+"_raw", article.Md, 1)
	c.Text(article.Md)
}

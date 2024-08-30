package handlers

import (
	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog/log"

	"github.com/pluja/blogo/frontend/templates"
	"github.com/pluja/blogo/internal/articles"
)

func HandleArticle(c iris.Context) {
	slug := c.Params().Get("slug")
	article, err := articles.GetFromFile(slug)
	if err != nil {
		log.Error().Err(err).Msg("error")
		c.RenderComponent(templates.ArticleNotFound(slug))
		return
	}

	c.RenderComponent(templates.Article(article))
}

func HandleRawArticle(c iris.Context) {
	slug := c.Params().Get("slug")
	article, err := articles.GetFromFile(slug)
	if err != nil {
		log.Error().Err(err).Msg("error")
		c.RenderComponent(templates.ArticleNotFound(slug))
		return
	}

	c.Text(article.Md)
}

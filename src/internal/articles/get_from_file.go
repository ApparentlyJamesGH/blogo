package articles

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"

	"github.com/pluja/blogo/internal/models"
	"github.com/pluja/blogo/internal/utils"
)

func GetFromFile(slug string) (models.Article, error) {
	articlesPath := viper.GetString("articlesPath")
	path := filepath.Join(articlesPath, slug+".md")

	var article models.Article

	// Read the markdown file
	content, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to read file")
		return article, err
	}

	var buf strings.Builder
	pContext := parser.NewContext()
	if err := markdown.Convert(content, &buf, parser.WithContext(pContext)); err != nil {
		return article, err
	}

	// Get the metadata fields
	metadata := meta.Get(pContext)

	// Handle drafts
	draftValue, exists := metadata["Draft"]
	if !exists {
		log.Warn().Msgf("%s has no draft value. Defaulting to false", path)
		draftValue = false
	}
	var draft bool
	switch articleDraft := draftValue.(type) {
	case bool:
		draft = articleDraft
	case string:
		if isDraft, err := strconv.ParseBool(articleDraft); err == nil && isDraft {
			draft = true
		} else if err != nil {
			log.Warn().Msgf("Could not parse draft value for %v, considering draft", path)
			draft = true
		}
	default:
		log.Err(err).Msgf("Could not parse draft value for %v, considering draft", path)
		draft = true
	}

	// Parse date
	dateString := utils.GetMapStringValue(metadata, "Date")
	date, err := time.Parse("2006-01-02 15:04", dateString)
	if err != nil {
		date, err = time.Parse("2006-01-02", dateString)
		if err != nil {
			log.Warn().Err(err).Msgf("Could not parse date for %v, using current time", path)
			date = time.Now()
		}
	}

	// Parse header image
	image := utils.GetMapStringValue(metadata, "Image")
	if image != "" && strings.HasPrefix(image, "/") {
		image = fmt.Sprintf("%v%v", viper.GetString("baseURL"), image)
	}

	// Fill article Data
	article = models.Article{
		Date:     date,
		Draft:    draft,
		Image:    image,
		Title:    utils.GetMapStringValue(metadata, "Title"),
		Author:   utils.GetMapStringValue(metadata, "Author"),
		Summary:  utils.GetMapStringValue(metadata, "Summary"),
		Layout:   utils.GetMapStringValue(metadata, "Layout"),
		NostrUrl: utils.GetMapStringValue(metadata, "NostrUrl"),
	}

	if tags, ok := metadata["Tags"].([]interface{}); ok {
		for _, tag := range tags {
			if strTag, ok := tag.(string); ok {
				article.Tags = append(article.Tags, strTag)
			} else {
				log.Warn().Msgf("Could not parse tag %v for %v", tag, path)
			}
		}
	}

	html, md, err := GetArticleContent(path)
	if err != nil {
		return models.Article{}, err
	}

	article.Html = html
	article.Md = md

	article.Slug = slug

	return article, nil
}

// Parses a .md file and returns the HTML and the raw markdown
func GetArticleContent(path string) (template.HTML, string, error) {
	md, err := os.ReadFile(path)
	if err != nil {
		return template.HTML(""), "", err
	}

	// Remove everything in the yaml metadata block

	var htmlBuf bytes.Buffer
	err = markdown.Convert(md, &htmlBuf)
	if err != nil {
		return template.HTML(""), "", err
	}
	html := htmlBuf.Bytes()
	return template.HTML(html), string(md), nil
}

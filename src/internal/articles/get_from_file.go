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

// GetFromFile retrieves an article from a file and parses its content.
func GetFromFile(slug string, parseContent bool) (models.Article, error) {
	articlesPath := viper.GetString("articles.path")
	path := filepath.Join(articlesPath, slug+".md")

	content, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to read file")
		return models.Article{}, fmt.Errorf("failed to read file: %w", err)
	}

	var buf strings.Builder
	pContext := parser.NewContext()
	if err := markdown.Convert(content, &buf, parser.WithContext(pContext)); err != nil {
		return models.Article{}, fmt.Errorf("failed to convert markdown: %w", err)
	}

	metadata := meta.Get(pContext)

	article := models.Article{
		Draft:   parseBoolField(metadata, "Draft", false),
		Nostr:   parseBoolField(metadata, "Nostr", true),
		Date:    parseDateField(metadata, path),
		Image:   parseImageField(metadata),
		Title:   utils.GetMapStringValue(metadata, "Title"),
		Author:  utils.GetMapStringValue(metadata, "Author"),
		Summary: utils.GetMapStringValue(metadata, "Summary"),
		Layout:  utils.GetMapStringValue(metadata, "Layout"),
		Tags:    parseTagsField(metadata, path),
		Slug:    slug,
	}

	if parseContent {
		html, md, err := GetArticleContent(path)
		if err != nil {
			return models.Article{}, fmt.Errorf("failed to get article content: %w", err)
		}
		article.Html = html
		if article.Image != "" {
			log.Debug().Msgf("Image: %s", article.Image)
			article.Html = template.HTML(fmt.Sprintf("\n<img src='%s' alt='header image'>\n", template.HTMLEscapeString(article.Image))) + article.Html
		}
		article.Md = md
	}

	return article, nil
}

// parseBoolField parses a boolean field from metadata.
func parseBoolField(metadata map[string]interface{}, key string, defaultNoExist bool) bool {
	value, exists := metadata[key]
	if !exists {
		log.Debug().Msgf("No %s value. Defaulting to %v", key, defaultNoExist)
		return defaultNoExist
	}

	switch v := value.(type) {
	case bool:
		return v
	case string:
		if parsed, err := strconv.ParseBool(v); err == nil {
			return parsed
		}
		log.Warn().Msgf("Could not parse %s value for %v, considering true", key, v)
	default:
		log.Warn().Msgf("Could not parse %s value for %v, considering true", key, v)
	}
	return true
}

// parseDateField parses the date field from metadata.
func parseDateField(metadata map[string]interface{}, path string) time.Time {
	dateString := utils.GetMapStringValue(metadata, "Date")
	for _, layout := range []string{"2006-01-02 15:04", "2006-01-02"} {
		if date, err := time.Parse(layout, dateString); err == nil {
			return date
		}
	}
	log.Warn().Msgf("Could not parse date for %v, using current time", path)
	return time.Now()
}

// parseImageField parses the image field from metadata.
func parseImageField(metadata map[string]interface{}) string {
	image := utils.GetMapStringValue(metadata, "Image")
	if image != "" && strings.HasPrefix(image, "/") {
		return fmt.Sprintf("%v%v", viper.GetString("host"), image)
	}
	return image
}

// parseTagsField parses the tags field from metadata.
func parseTagsField(metadata map[string]interface{}, path string) []string {
	var tags []string
	if tagList, ok := metadata["Tags"].([]interface{}); ok {
		for _, tag := range tagList {
			if strTag, ok := tag.(string); ok {
				tags = append(tags, strTag)
			} else {
				log.Warn().Msgf("Could not parse tag %v for %v", tag, path)
			}
		}
	}
	return tags
}

// GetArticleContent parses a .md file and returns the HTML and the raw markdown.
func GetArticleContent(path string) (template.HTML, string, error) {
	md, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := markdown.Convert(md, &htmlBuf); err != nil {
		return "", "", fmt.Errorf("failed to convert markdown: %w", err)
	}

	return template.HTML(htmlBuf.String()), string(md), nil
}

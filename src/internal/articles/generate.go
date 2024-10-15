// internal/articles/generate.go
package articles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GenerateNewArticle(slug, title, author, tags, summary string) {
	articlesPath := viper.GetString("articles.path")
	fileName := fmt.Sprintf("%s.md", slug)
	filePath := filepath.Join(articlesPath, fileName)

	if title == "" {
		title = slugToTitle(slug)
	}

	content := generateArticleContent(title, author, tags, summary)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating article: %v\n", err)
		return
	}

	fmt.Printf("Article created successfully: %s\n", filePath)
}

func slugToTitle(slug string) string {
	words := strings.Split(slug, "-")
	for i, word := range words {
		words[i] = cases.Title(language.Und).String(word)
	}
	return strings.Join(words, " ")
}

func generateArticleContent(title, author, tags, summary string) string {
	currentTime := time.Now().Format("2006-01-02 15:04")
	tagList := strings.Split(tags, ",")
	for i, tag := range tagList {
		tagList[i] = strings.TrimSpace(tag)
	}

	content := fmt.Sprintf(`---
Author: %s
Date: %s
Draft: true
Layout: post
Image: 
Summary: %s
Tags:
%s
Title: %s
---

# A sample post

This is a sample blog post. Edit it with your own content!
`, author, currentTime, summary, formatTags(tagList), title)

	return content
}

func formatTags(tags []string) string {
	var formattedTags string
	for _, tag := range tags {
		if tag != "" {
			formattedTags += fmt.Sprintf("- %s\n", tag)
		}
	}
	return formattedTags
}

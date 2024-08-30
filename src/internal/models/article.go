package models

import (
	"html/template"
	"time"
)

type Article struct {
	Title    string
	Author   string
	Summary  string
	Tags     []string
	Image    string
	Date     time.Time
	Slug     string
	Draft    bool
	Layout   string
	Md       string
	Html     template.HTML
	NostrUrl string
}

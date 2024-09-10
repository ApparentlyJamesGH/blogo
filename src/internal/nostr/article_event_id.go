package nostr

import (
	"crypto/md5"
	"fmt"

	"github.com/pluja/blogo/internal/models"
)

// md5 hash the title and slug to get a unique ID
func GetEventId(article models.Article) string {
	id := fmt.Sprintf("%x", md5.Sum([]byte(
		nostrPk+article.Title+article.Author,
	)))
	return id
}

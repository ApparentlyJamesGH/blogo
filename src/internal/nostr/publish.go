package nostr

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/pluja/blogo/internal/models"
	"github.com/pluja/blogo/internal/utils"
)

var NostrPublications = make(map[string]string) // Initialize the map

// Publishes the article to Nostr if enabled and not yet published
func PublishArticle(article models.Article) error {
	if article.Slug == "about" {
		log.Info().Msg("Not publishing about page to Nostr")
		return nil
	}

	if os.Getenv("PUBLISH_TO_NOSTR") == "false" {
		log.Info().Msg("PUBLISH_TO_NOSTR is set to false. Not publishing...")
		return nil
	}

	// If the Nostr field is set to something that evaluates as false, we don't publish
	if !article.Nostr {
		log.Info().Msgf("Nostr value (%v) is set to False. Not publishing...", article.Nostr)
		return nil
	}

	// If the article is a draft we don't publish
	if !article.Draft {
		log.Info().Msgf("Publishing %q to Nostr", article.Slug)
		naddr, err := nostrPublish(article)
		if err != nil {
			log.Err(err).Msg("Could not publish to Nostr")
		} else {
			NostrPublications[article.Slug] = naddr
		}
	} else {
		log.Printf("Won't publisht this to Nostr: it's a draft")
	}
	return nil
}

// Publishes an article of type models.Article to Nostr.
func nostrPublish(ad models.Article) (string, error) {
	// Wipe the YAML Metadata block from the article
	sections := strings.SplitN(string(ad.Md), "---", 3)
	if len(sections) >= 3 {
		ad.Md = sections[2]
	}

	// Add the article original URL to the top of the article
	ad.Md = fmt.Sprintf("> [Read the original blog post](%v)\n\n", utils.CreateURL(viper.GetString("base_url"), "/p/", ad.Slug)) + ad.Md

	id := GetEventId(ad)
	// Create the Nostr event
	tags := nostr.Tags{
		nostr.Tag{"d", id},
		nostr.Tag{"title", ad.Title},
		nostr.Tag{"summary", ad.Summary},
		nostr.Tag{"image", ad.Image},
		nostr.Tag{"slug", ad.Slug},
		nostr.Tag{"client", "blogo"},
		nostr.Tag{"published_at", strconv.FormatInt(ad.Date.Unix(), 10)},
	}

	articleTags := nostr.Tags{}
	for _, tag := range ad.Tags {
		articleTags = append(articleTags, nostr.Tag{"t", tag})
	}
	tags = append(tags, articleTags...)

	ev := nostr.Event{
		PubKey:    nostrPk,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindArticle,
		Tags:      tags,
		Content:   ad.Md,
	}

	// Sign the event
	err := ev.Sign(nostrSk)
	if err != nil {
		log.Err(err).Msg("Could not sign event")
		return "", err
	}

	// Publish the event to the relays
	ctx := context.Background()
	connected := false
	published := false
	if os.Getenv("DEV") == "true" {
		// In development mode, mock the Nostr publish
		connected = true
		published = true
	} else {
		for _, url := range viper.GetStringSlice("nostr.relays") {
			relay, err := nostr.RelayConnect(ctx, url)
			if err != nil {
				log.Err(err).Msgf("failed to connect to relay %v:", url)
				continue
			}
			connected = true
			if err := relay.Publish(ctx, ev); err != nil {
				log.Warn().Err(err).Msgf("failed to publish to %v", url)
				continue
			}

			published = true
			fmt.Printf("published %q to %s\n", ev.ID, url)
		}
	}

	// Return an error if we were unable to publish to any relay
	if !connected || !published {
		return "", fmt.Errorf("unable to publish %q to Nostr, connected: %v, published: %v", ev.ID, connected, published)
	}

	// Encode the note ID to naddr format
	naddr, err := nip19.EncodeEntity(ev.PubKey, nostr.KindArticle, id, []string{})
	if err != nil {
		log.Err(err).Msg("Could not encode note ID")
		return ev.ID, err
	}
	return naddr, nil
}

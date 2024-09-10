package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/pluja/blogo/internal/articles"
	"github.com/pluja/blogo/internal/cache"
	"github.com/pluja/blogo/internal/nostr"
	"github.com/pluja/blogo/internal/server"
)

var devMode *bool

func init() {
	godotenv.Load()
	initFlags()
	initLogger()
	loadConfig()
	articles.InitGoldmark()
	cache.Init()
	nostr.Init()
}

func main() {
	go articles.WatchArticles()

	if err := server.StartServer(); err != nil {
		log.Fatal().Err(err)
	}
}

func initLogger() {
	if *devMode {
		os.Setenv("DEV_MODE", "true")
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Printf("DEV_MODE: %v", os.Getenv("DEV_MODE"))
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func initFlags() {
	devMode = flag.Bool("dev", false, "activate dev mode")
	flag.Parse()
}

func loadConfig() {
	viper.SetConfigName("blogo")
	viper.SetConfigType("yaml")

	// Config paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("/blogo/")
	viper.AddConfigPath("$HOME/.blogo")

	// Env variables
	viper.SetEnvPrefix("BLOGO")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("powered_by_footer", true)
	viper.SetDefault("title", "Blogo")
	viper.SetDefault("description", "Welcome to my blogo")
	viper.SetDefault("timezone", "UTC")
	viper.SetDefault("theme", "blogo")
	viper.SetDefault("nostr.publish", false)
	viper.SetDefault("nostr.relays", []string{"wss://nostr-pub.wellorder.net", "wss://relay.damus.io", "wss://relay.nostr.band"})
	viper.SetDefault("articles_path", "/blogo/articles")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	viper.WatchConfig()

	log.Printf("%s", viper.GetString("title"))
}

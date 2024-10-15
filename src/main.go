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
	initLogger(false)
	loadConfig()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: blogo <command> [arguments]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		devMode := serveCmd.Bool("dev", false, "activate dev mode")
		serveCmd.Parse(os.Args[2:])
		serve(*devMode)
	case "new":
		newCmd := flag.NewFlagSet("new", flag.ExitOnError)
		title := newCmd.String("title", "", "Article title")
		author := newCmd.String("author", "", "Article author")
		tags := newCmd.String("tags", "", "Comma-separated list of tags")
		summary := newCmd.String("summary", "", "Article summary")

		newCmd.Parse(os.Args[2:])

		if newCmd.NArg() < 1 {
			fmt.Println("Usage: blogo new <slug> [flags]")
			newCmd.PrintDefaults()
			os.Exit(1)
		}

		slug := newCmd.Arg(0)
		if !articles.ValidateSlug(slug) {
			fmt.Println("Invalid slug. Use only alphanumeric characters and hyphens.")
			os.Exit(1)
		}

		articles.GenerateNewArticle(slug, *title, *author, *tags, *summary)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func serve(devMode bool) {
	initLogger(devMode)
	articles.InitGoldmark()
	cache.Init()
	nostr.Init()
	go articles.WatchArticles()

	if err := server.StartServer(); err != nil {
		log.Fatal().Err(err)
	}
}

func initLogger(devMode bool) {
	if devMode {
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
	viper.SetDefault("title", "Anon Blog")
	viper.SetDefault("description", fmt.Sprintf("Welcome to %s's blog, enjoy!", viper.GetString("title")))
	viper.SetDefault("timezone", "UTC")
	viper.SetDefault("theme", "blogo")
	viper.SetDefault("nostr.publish", false)
	viper.SetDefault("nostr.relays", []string{"wss://nostr-pub.wellorder.net", "wss://relay.damus.io", "wss://relay.nostr.band"})
	viper.SetDefault("articles.path", "/blogo/articles")

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

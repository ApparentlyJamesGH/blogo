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
	"github.com/pluja/blogo/internal/server"
)

var devMode *bool

func init() {
	godotenv.Load()
	initFlags()
	initLogger()
	loadConfig()
	articles.InitGoldmark()
	cache.InitCache()
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
	viper.AddConfigPath(".")
	viper.AddConfigPath("/blogo/")
	viper.AddConfigPath("$HOME/.blogo")
	viper.SetEnvPrefix("BLOGO")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	log.Printf("%s", viper.GetString("title"))
}

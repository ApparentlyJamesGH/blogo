package server

import (
	"path"

	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog/log"

	"github.com/pluja/blogo/frontend/templates"
	"github.com/pluja/blogo/internal/server/handlers"
	"github.com/pluja/blogo/internal/utils"
)

func StartServer() error {
	r := iris.New()

	r.Use(iris.Compression)
	//r.Logger().SetLevel("debug")
	// r.Use(iris.Cache304(24 * 60 * 60))

	serverAddress := utils.Getenv("SERVER_ADDRESS", ":1337")

	// Static routes
	rootDir := utils.Getenv("ROOT_DIR", "./")
	r.Favicon(path.Join(rootDir, "/frontend/static", "/assets/favicon.webp"))
	r.HandleDir("/static", iris.Dir(path.Join(rootDir, "/frontend/static")))

	// UI Handlers
	r.Get("/", iris.Component(templates.Index()))
	r.Get("/blog", iris.Component(templates.Blog()))
	r.Get("/rss", handlers.HandleFeed)
	r.Get("/atom", handlers.HandleFeed)
	r.Get("/json", handlers.HandleFeed)
	r.Get("/p/{slug}", handlers.HandleArticle)
	r.Get("/p/{slug}/raw", handlers.HandleRawArticle)

	log.Info().Msgf("Starting server at http://127.0.0.1%s", serverAddress)
	return r.Listen(serverAddress)
}

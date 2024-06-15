package handler

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/h24s_24/repository"
)

func Start() {
	e := echo.New()

	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("failed to get db: %v\n", err)
	}

	pr := repository.NewPostRepository(db)
	if err != nil {
		log.Fatalf("failed to get post repository: %v\n", err)
	}
	rr := repository.NewReactionRepository(db)
	if err != nil {
		log.Fatalf("failed to get reaction repository: %v\n", err)
	}

	ph := &PostHandler{PostRepository: pr, ReactionRepository: rr}

	e.GET("/health", func(c echo.Context) error { return c.String(200, "OK") })

	api := e.Group("/api")
	api.POST("/posts", ph.PostPostsHandler)
	api.GET("/posts", ph.GetPostsHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

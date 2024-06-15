package handler

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	rh := &ReactionHandler{rr: rr}

	e.Use(middleware.Logger(), middleware.Recover())

	e.GET("/health", func(c echo.Context) error { return c.String(200, "OK") })

	api := e.Group("/api")
	api.POST("/posts", ph.PostPostsHandler)
	api.GET("/posts", ph.GetPostsHandler)

	api.POST("/posts/:postID/reactions/:reactionID", rh.PostReactionHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

var errNoUsername = errors.New("no username")

func getUsername(c echo.Context) (string, error) {
	username := c.Request().Header.Get("X-Forwarded-User")
	if username != "" {
		return username, nil
	}

	local, err := strconv.ParseBool(os.Getenv("LOCAL"))
	if err != nil || !local {
		return "", errNoUsername
	}
	return "testuser", nil
}

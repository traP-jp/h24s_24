package handler

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traP-jp/h24s_24/converter"
	"github.com/traP-jp/h24s_24/converter/mock"
	"github.com/traP-jp/h24s_24/repository"
)

func Start() {
	e := echo.New()

	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("failed to get db: %v\n", err)
	}

	pr := repository.NewPostRepository(db)
	rr := repository.NewReactionRepository(db)
	ur := repository.NewUserRepository(db)

	// ローカルのときはモックを使う
	var cvt PostConverter
	if local, err := strconv.ParseBool(os.Getenv("LOCAL")); err == nil && local {
		log.Println("using mock converter")
		cvt = &mock.MockConverter{}
	} else {
		cvt, err = converter.NewOpenAI()
		if err != nil {
			log.Fatalf("failed to get OpenAI converter: %v\n", err)
		}
	}

	ph := &PostHandler{PostRepository: pr, ReactionRepository: rr, pc: cvt}
	rh := &ReactionHandler{rr: rr}
	th := &TrendHandler{rr: rr, pr: pr}
	uh := &UserHandler{rr: rr, ur: ur}

	e.Use(middleware.Logger(), middleware.Recover())

	e.GET("/health", func(c echo.Context) error { return c.String(200, "OK") })

	api := e.Group("/api")
	api.Use(userNameMiddleware)
	api.POST("/posts", ph.PostPostsHandler)
	api.GET("/posts", ph.GetPostsHandler)
	api.GET("/posts/:postID", ph.GetPostHandler)

	api.POST("/posts/:postID/reactions/:reactionID", rh.PostReactionHandler)

	api.GET("/trend", th.GetTrendHandler)

	api.GET("/users/:userName", uh.GetUserHandler)
	api.GET("/me", uh.GetMeHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

var errNoUsername = errors.New("no username")

func getUserName(c echo.Context) (string, error) {
	userName, ok := c.Get(userNameCtxKey).(string)
	if !ok {
		return "", errNoUsername
	}

	return userName, nil
}

package handler

import (
	"cmp"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TrendHandler struct {
	pr PostRepository
	rr ReactionRepository
}

type postReactionsCount struct {
	PostID uuid.UUID
	Count  int
}

type reaction struct {
	ReactionID int `json:"id"`
	Count      int `json:"count"`
}

type getTrendResponse struct {
	PostID           uuid.UUID   `json:"id"`
	Username         string      `json:"user_name"`
	OriginalMessage  string      `json:"original_message"`
	ConvertedMessage string      `json:"converted_message"`
	Reactions        []*reaction `json:"reactions"`
	RootID           uuid.UUID   `json:"root_id"`
}

func (tr *TrendHandler) GetTrendHandler(c echo.Context) error {
	since := time.Now()
	timeSpan := 3600 * time.Second
	reactionLimit := 250
	postLimit := 30

	ctx := c.Request().Context()
	countsMap, err := tr.rr.GetReactionsCount(ctx, since, timeSpan, reactionLimit)
	if err != nil {
		log.Printf("failed to get scores: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get scores")
	}

	sortedPostCounts := make([]*postReactionsCount, 0, len(countsMap))

	for postID, reactionsCount := range countsMap {
		sortedPostCounts = append(sortedPostCounts, &postReactionsCount{postID, reactionsCount})
	}
	slices.SortFunc(sortedPostCounts, func(a, b *postReactionsCount) int {
		n := cmp.Compare(b.Count, a.Count)
		if n != 0 {
			return n
		}
		return cmp.Compare(a.PostID.ID(), b.PostID.ID())
	})

	posts := make([]*getTrendResponse, min(postLimit, len(sortedPostCounts)))
	for i, v := range sortedPostCounts {
		post, err := tr.pr.GetPostByID(ctx, v.PostID)
		if err != nil {
			log.Printf("failed to get post: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
		}

		reactions, err := tr.rr.GetReactionsByPostID(ctx, v.PostID)
		if err != nil {
			log.Printf("failed to get reactions: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get reactions")
		}

		reactionsMap := make(map[int]int)
		for _, r := range reactions {
			reactionsMap[r.ReactionID] += r.Count
		}
		reactionsSlice := make([]*reaction, 0, len(reactionsMap))
		for reactionID, count := range reactionsMap {
			reactionsSlice = append(reactionsSlice, &reaction{reactionID, count})
		}

		posts[i] = &getTrendResponse{
			PostID:           post.ID,
			Username:         post.UserName,
			OriginalMessage:  post.OriginalMessage,
			ConvertedMessage: post.ConvertedMessage,
			Reactions:        reactionsSlice,
			RootID:           post.RootID,
		}
	}
	return c.JSON(http.StatusOK, posts)
}

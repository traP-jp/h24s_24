package random

import (
	"context"
	"fmt"

	"github.com/traP-jp/h24s_24/domain"
)

type PostRepository interface {
	GetRandomPost(ctx context.Context) (domain.Post, error)
}

type Converter struct {
	pr PostRepository
}

func NewConverter(pr PostRepository) *Converter {
	return &Converter{pr: pr}
}

func (cvt *Converter) ConvertMessage(ctx context.Context, originalMessage string) (string, error) {
	post, err := cvt.pr.GetRandomPost(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get random post: %w", err)
	}

	converted := fmt.Sprintf("%s\nそれはそうと、%s", originalMessage, post.ConvertedMessage)

	return converted, nil
}

package mock

import "context"

type MockConverter struct{}

func (mc *MockConverter) ConvertMessage(ctx context.Context, originalMessage string) (string, error) {
	return originalMessage + " (converted by mock)", nil
}

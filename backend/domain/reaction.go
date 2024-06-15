package domain

type Reaction struct {
	PostID int
	Count  int
	Users  []string
}

package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Status  int32
	Tags    []string
}

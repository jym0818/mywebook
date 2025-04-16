package domain

type Article struct {
	Content string
	Title   string
	Author  Author
}

type Author struct {
	Id   int64
	Name string
}

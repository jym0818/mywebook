package domain

type Article struct {
	Id      int64
	Content string
	Title   string
	Author  Author
}

type Author struct {
	Id   int64
	Name string
}

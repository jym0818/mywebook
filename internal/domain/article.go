package domain

type Article struct {
	Id      int64
	Content string
	Title   string
	Author  Author
	Status  ArticleStatus
}

type Author struct {
	Id   int64
	Name string
}
type ArticleStatus uint8

const (
	ArticleStatusUnKnown ArticleStatus = iota
	//未发表
	ArticleStatusUnPublished
	//已发表
	ArticleStatusPublished
	//仅自己可见
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

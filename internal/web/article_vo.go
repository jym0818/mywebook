package web

import "github.com/jym/mywebook/internal/domain"

type CollectReq struct {
	Id  int64 `json:"id"`
	Cid int64 `json:"cid"`
}
type LikeReq struct {
	Id   int64 `json:"id"`
	Like bool  `json:"like"`
}
type ArticleVO struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	// 内容
	Content string `json:"content"`
	// 注意一点，状态这个东西，可以是前端来处理，也可以是后端处理
	// 0 -> unknown -> 未知状态
	// 1 -> 未发表，手机 APP 这种涉及到发版的问题，那么后端来处理
	// 涉及到国际化，也是后端来处理
	Status     uint8  `json:"status"`
	AuthorId   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	Ctime      string `json:"ctime"`
	Utime      string `json:"utime"`

	// 点赞之类的信息
	LikeCnt    int64 `json:"likeCnt"`
	CollectCnt int64 `json:"collectCnt"`
	ReadCnt    int64 `json:"readCnt"`

	// 个人是否点赞的信息
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type ListReq struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  uint8  `json:"status"`
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
		Status: domain.ArticleStatus(req.Status),
	}
}

package service

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"github.com/jym/mywebook/internal/domain"
	"github.com/jym/mywebook/internal/repository"
	service2 "github.com/jym/webook-interactive/service"
	"math"
	"time"
)

type RankingService interface {
	//设置热榜
	TopN(ctx context.Context) error
}

type BatchRankingService struct {
	artSvc    ArticleService
	intrSvc   service2.InteractiveService
	batchSize int
	n         int
	//规避计算出负值 否则一开始没有点赞数据都是负数了进不去队列了
	scoreFunc func(t time.Time, likeCnt int64) float64
	repo      repository.RankingRepository
}

func NewBatchRankingService(artSvc ArticleService, intrSvc service2.InteractiveService) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		intrSvc:   intrSvc,
		batchSize: 100,
		n:         100,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			ms := time.Since(t).Seconds()
			return float64(likeCnt+1) / math.Pow(ms+2, 1.5)
		},
	}
}

// 有返回值才方便测试
func (svc *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	//我们只取7天内的数据大于7天的数据不是热榜
	now := time.Now()

	offset := 0

	type Score struct {
		art   domain.Article
		score float64
	}

	//小顶堆实现的队列
	//非并发安全
	topN := queue.NewConcurrentPriorityQueue[Score](svc.n, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score == dst.score {
			return 0
		} else {
			return -1
		}
	})
	for {
		//拿一批数据 只取7天内的数据
		arts, err := svc.artSvc.ListPub(ctx, now, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts, func(idx int, src domain.Article) int64 {
			return src.Id
		})
		//再拿到这批数据的点赞数
		intrs, err := svc.intrSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}
		//合并结算score
		for _, art := range arts {
			intr, _ := intrs[art.Id]

			score := svc.scoreFunc(art.Utime, intr.LikeCnt)
			//排序
			//我要考虑score在不在前100名
			//1.先直接入队
			err = topN.Enqueue(Score{art, score})
			//队列满了 说明需要淘汰了
			if err == queue.ErrOutOfCapacity {
				//出列----也就是最小值
				val, _ := topN.Dequeue()
				//判断 如果比最小值大就入列
				if val.score < score {
					err = topN.Enqueue(Score{art, score})
				} else {
					//不符合 再放回去
					err = topN.Enqueue(val)
				}
			}
		}

		//一批数据处理完了，如何判断是不是要进入下一批，怎么判断是否还有数据
		//分页查询结束的判断方式  如果查出的数量少于batchSize说明这是最后一页了
		if len(arts) < svc.batchSize || now.Sub(arts[len(arts)-1].Utime).Hours() > 7*24 {
			//只取7天内数据
			break
		}
		//否则更新offset
		offset = offset + len(arts)

	}
	//得到结果
	res := make([]domain.Article, 0, svc.n)
	for i := svc.n - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			//说明取完了不够100
			break
		}
		res[i] = val.art

	}
	return res, nil
}
func (svc *BatchRankingService) TopN(ctx context.Context) error {
	arts, err := svc.topN(ctx)
	if err != nil {
		return err
	}
	//放当缓存里面  这个项目中我们不会放在数据库内

	return svc.repo.ReplaceTopN(ctx, arts)
}

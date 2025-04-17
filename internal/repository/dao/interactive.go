package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrNotFound = gorm.ErrRecordNotFound

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error)
	GetCollectionInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error)
	Get(ctx context.Context, biz string, id int64) (Interactive, error)
}
type GORMInteractiveDAO struct {
	db *gorm.DB
}

func (dao *GORMInteractiveDAO) Get(ctx context.Context, biz string, id int64) (Interactive, error) {
	var res Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ?", biz, id).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetCollectionInfo(ctx context.Context, biz string, id int64, uid int64) (UserCollectionBiz, error) {
	var user UserCollectionBiz
	err := dao.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ?", biz, id, uid).First(&user).Error
	return user, err
}

func (dao *GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
	var user UserLikeBiz
	err := dao.db.WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ?", biz, bizId, uid).First(&user).Error
	return user, err
}

func (dao *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.Utime = now
	cb.Ctime = now
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := dao.db.WithContext(ctx).Create(&cb).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt`+1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
			Biz:        cb.Biz,
			BizId:      cb.BizId,
		}).Error
	})
}

func (dao *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).Where("biz =? AND biz_id = ? AND uid = ?", biz, id, uid).
			Updates(map[string]interface{}{
				"status": 0,
				"utime":  now,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).Where("biz = ? AND biz_id = ?", biz, id).Updates(map[string]interface{}{
			"utime":    now,
			"like_cnt": gorm.Expr("like_cnt - 1"),
		}).Error
	})
	return err
}

func (dao *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"status": 1,
				"utime":  now,
			}),
		}).Create(&UserLikeBiz{
			Uid:    uid,
			Ctime:  now,
			Utime:  now,
			Biz:    biz,
			BizId:  id,
			Status: 1,
		}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("`like_cnt`+1"),
				"utime":    now,
			}),
		}).Create(&Interactive{
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
			Biz:     biz,
			BizId:   id,
		}).Error
	})
	return err
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{
		db: db,
	}
}
func (dao *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	var intr Interactive
	intr.BizId = bizId
	intr.Biz = biz
	intr.Utime = now
	intr.Ctime = now
	intr.ReadCnt = 1
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"utime":    now,
			"read_cnt": gorm.Expr("`read_cnt`+1"),
		}),
	}).Create(&intr).Error
}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz        string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`
	ReadCnt    int64
	CollectCnt int64
	LikeCnt    int64
	Ctime      int64
	Utime      int64
}
type UserLikeBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 三个构成唯一索引
	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	Uid   int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	// 依旧是只在 DB 层面生效的状态
	// 1- 有效，0-无效。软删除的用法
	Status uint8
	Ctime  int64
	Utime  int64
}
type Collection struct {
	Id   int64  `gorm:"primaryKey,autoIncrement"`
	Name string `gorm:"type=varchar(1024)"`
	Uid  int64  `gorm:""`

	Ctime int64
	Utime int64
}

type UserCollectionBiz struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 收藏夹 ID
	// 作为关联关系中的外键，我们这里需要索引
	Cid   int64  `gorm:"index"`
	BizId int64  `gorm:"uniqueIndex:biz_type_id_uid"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id_uid"`
	// 这算是一个冗余，因为正常来说，
	// 只需要在 Collection 中维持住 Uid 就可以
	Uid   int64 `gorm:"uniqueIndex:biz_type_id_uid"`
	Ctime int64
	Utime int64
}

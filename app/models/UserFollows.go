package models

import (
	"context"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/libs/logger"
)

// 用户关注
type UserFollows struct {
	ID          int `json:"id" gorm:"column:id"`
	Uid         int `json:"uid" gorm:"column:uid"`                   // uid
	ContentType int `json:"content_type" gorm:"column:content_type"` // 1文章，2回答
	RelationId  int `json:"relation_id" gorm:"column:relation_id"`   // 关联id
}

func NewUserFollows() *UserFollows {
	return new(UserFollows)
}
func (*UserFollows) TableName() string {
	return "user_follows"
}
func (a *UserFollows) Insert(ctx context.Context) {
	err := core.DB.Create(&a).Error
	if err != nil {
		logger.Info(ctx, "db insert fail", err)
	}
}
func (a *UserFollows) Save(ctx context.Context) {
	core.DB.Save(a)
}

//获取分页列表
func (a *UserFollows) PageList(where map[string]interface{}, page int, pagesize int, order string) (result []UserFollows) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	core.DB.Where(where).Order(order).Offset(offset).Limit(pagesize).Find(&result)
	return
}
func (a *UserFollows) Del(ctx context.Context) {
	core.DB.Delete(a)
}

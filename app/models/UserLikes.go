package models

import (
	"context"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/libs/logger"
)

// 用户点赞
type UserLikes struct {
	ID          int `json:"id" gorm:"column:id"`
	Uid         int `json:"uid" gorm:"column:uid"`                   // uid
	ContentType int `json:"content_type" gorm:"column:content_type"` // 1文章，2回答
	RelationId  int `json:"relation_id" gorm:"column:relation_id"`   // 关联id
}

func NewUserLikes() *UserLikes {
	return new(UserLikes)
}
func (*UserLikes) TableName() string {
	return "user_likes"
}
func (a *UserLikes) Insert(ctx context.Context) {
	err := core.DB.Create(&a).Error
	if err != nil {
		logger.Info(ctx, "db insert fail", err)
	}
}
func (a *UserLikes) Save(ctx context.Context) {
	core.DB.Save(a)
}

//获取分页列表
func (a *UserLikes) PageList(where map[string]interface{}, page int, pagesize int, order string) (result []UserLikes) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	core.DB.Where(where).Order(order).Offset(offset).Limit(pagesize).Find(&result)
	return
}
func (a *UserLikes) Del(ctx context.Context) {
	core.DB.Delete(a)
}

package models

import (
	"context"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/libs/logger"
)

// 标签关联表
type TagRelation struct {
	ID          int `json:"id" gorm:"column:id"`
	TagId       int `json:"tag_id" gorm:"column:tag_id"`             // tagID
	ContentType int `json:"content_type" gorm:"column:content_type"` // 标签类型1行业，2分类
	RelationId  int `json:"relation_id" gorm:"column:relation_id"`   // 关联id
}

func NewTagRelation() *TagRelation {
	return new(TagRelation)
}
func (*TagRelation) TableName() string {
	return "tag_relation"
}
func (a *TagRelation) Insert(ctx context.Context) {
	err := core.DB.Create(&a).Error
	if err != nil {
		logger.Info(ctx, "db insert fail", err)
	}
}
func (a *TagRelation) Save(ctx context.Context) {
	core.DB.Save(a)
}

//获取分页列表
func (a *TagRelation) PageList(where map[string]interface{}, page int, pagesize int, order string) (result []TagRelation) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	core.DB.Where(where).Order(order).Offset(offset).Limit(pagesize).Find(&result)
	return
}
func (a *TagRelation) Del(ctx context.Context) {
	core.DB.Delete(a)
}

package models

import (
	"context"
	"time"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/libs/logger"
)

// 问题
type Question struct {
	ID         int       `json:"id" gorm:"column:id"`
	UniqId     string    `json:"uniq_id" gorm:"column:uniq_id"`                      // 唯一ID
	Uid        int       `json:"uid" gorm:"column:uid"`                              // 用户ID
	Title      string    `json:"title" gorm:"column:title"`                          // 文章标题
	Desc       string    `json:"desc" gorm:"column:desc"`                            // 简述
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;default:null"` // 创建时间
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;default:null"` // 更新时间
	LikeNums   int       `json:"like_nums" gorm:"column:like_nums"`                  // 点赞数
	ViewNums   int       `json:"view_nums" gorm:"column:view_nums"`                  // 浏览量
	FollowNums int       `json:"follow_nums" gorm:"column:follow_nums"`              // 关注数
	AnswerNums int       `json:"answer_nums" gorm:"column:answer_nums"`              // 回答数
	State      int       `json:"state" gorm:"column:state"`                          // 状态 0为禁用1为启用
}

func NewQuestion() *Question {
	return new(Question)
}
func (*Question) TableName() string {
	return "question"
}
func (a *Question) Insert(ctx context.Context) {
	err := core.DB.Create(&a).Error
	if err != nil {
		logger.Info(ctx, "db insert fail", err)
	}
}
func (a *Question) Save(ctx context.Context) {
	core.DB.Save(a)
}

//唯一id获取记录
func (a *Question) GetByUniqID(ctx context.Context, uniqID string) {
	core.DB.Where("uniq_id = ?", uniqID).First(a)
}

//获取分页列表
func (a *Question) PageList(where map[string]interface{}, page int, pagesize int, order string) (result []Question) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	core.DB.Where(where).Order(order).Offset(offset).Limit(pagesize).Find(&result)
	return
}
func (a *Question) Del(ctx context.Context) {
	core.DB.Delete(a)
}

package models

import (
	"context"
	"time"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/libs/logger"
)

// 评论
type Discuss struct {
	ID          int       `json:"id" gorm:"column:id"`
	Uid         int       `json:"uid" gorm:"column:uid"`                              // 用户ID
	ParentType  int       `json:"parent_type" gorm:"column:parent_type"`              // 评论类型 1文章，2回答
	ParentId    int       `json:"parent_id" gorm:"column:parent_id"`                  // 文章或回答id
	Content     string    `json:"content" gorm:"column:content"`                      // 正文
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;default:null"` // 创建时间
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;default:null"` // 更新时间
	LikeNums    int       `json:"like_nums" gorm:"column:like_nums"`                  // 点赞数
	FollowNums  int       `json:"follow_nums" gorm:"column:follow_nums"`              // 关注数
	CommentNums int       `json:"comment_nums" gorm:"column:comment_nums"`            // 回答数
	State       int       `json:"state" gorm:"column:state"`                          // 状态 0为禁用1为启用
}

func NewDiscuss() *Discuss {
	return new(Discuss)
}
func (*Discuss) TableName() string {
	return "discuss"
}
func (a *Discuss) Insert(ctx context.Context) {
	err := core.DB.Create(&a).Error
	if err != nil {
		logger.Info(ctx, "db insert fail", err)
	}
}
func (a *Discuss) Save(ctx context.Context) {
	core.DB.Save(a)
}

//获取分页列表
func (a *Discuss) PageList(where map[string]interface{}, page int, pagesize int, order string) (result []Discuss) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	core.DB.Where(where).Order(order).Offset(offset).Limit(pagesize).Find(&result)
	return
}
func (a *Discuss) Del(ctx context.Context) {
	core.DB.Delete(a)
}

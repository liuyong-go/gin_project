package models

import (
	"context"
	"time"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/libs/logger"
	"gorm.io/gorm"
)

// 文章
type Article struct {
	ID          int       `json:"id"`
	UniqId      string    `json:"uniq_id"`                                            // 唯一ID
	Uid         int       `json:"uid"`                                                // 用户ID
	Title       string    `json:"title"`                                              // 文章标题
	Desc        string    `json:"desc"`                                               // 简述
	Content     string    `json:"content"`                                            // 正文
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time;default:null"` // 创建时间
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time;default:null"` // 更新时间
	LikeNums    int       `json:"like_nums"`                                          // 点赞数
	ViewNums    int       `json:"view_nums"`                                          // 浏览量
	FollowNums  int       `json:"follow_nums"`                                        // 关注数
	CommentNums int       `json:"comment_nums"`                                       // 评论数
	State       int       `json:"state"`                                              // 状态 0为禁用1为启用
}

func NewArticle() *Article {
	return new(Article)
}
func (*Article) TableName() string {
	return "article"
}
func (a *Article) Insert(ctx context.Context) {
	err := core.DB.Create(&a).Error
	if err != nil {
		logger.Info(ctx, "db insert fail", err)
	}
}
func (a *Article) Save(ctx context.Context) {
	core.DB.Save(a)
}

//字段值加一
func (a *Article) Incr(ctx context.Context, field string) {
	core.DB.Model(a).UpdateColumn(field, gorm.Expr(field+" + ?", 1))
}

//唯一id获取记录
func (a *Article) GetByUniqID(ctx context.Context, uniqID string) {
	core.DB.Where("uniq_id = ?", uniqID).First(a)
}

//获取分页列表
func (a *Article) PageList(where map[string]interface{}, page int, pagesize int, order string) (result []Article) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	core.DB.Where(where).Order(order).Offset(offset).Limit(pagesize).Find(&result)
	return
}
func (a *Article) Del(ctx context.Context) {
	core.DB.Delete(a)
}

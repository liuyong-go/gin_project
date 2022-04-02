package models

import (
	"context"
	"time"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/libs/logger"
)

// 用户微信表
type User struct {
	ID           int       `json:"id" gorm:"column:id"`                                // 主键
	Openid       string    `json:"openid" gorm:"column:openid"`                        // openid
	Unionid      string    `json:"unionid" gorm:"column:unionid"`                      // unionid
	Nickname     string    `json:"nickname" gorm:"column:nickname"`                    // 昵称
	Avatar       string    `json:"avatar" gorm:"column:avatar"`                        // 头像
	Gender       int       `json:"gender" gorm:"column:gender"`                        // 性别 0：未知、1：男、2：女
	City         string    `json:"city" gorm:"column:city"`                            // 市
	Province     string    `json:"province" gorm:"column:province"`                    // 省
	Country      string    `json:"country" gorm:"column:country"`                      // 国家
	Introduction string    `json:"introduction" gorm:"column:introduction"`            // 自我介绍
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time;default:null"` // 注册时间
	UpdateTime   time.Time `json:"update_time" gorm:"column:update_time;default:null"` // 更新时间
	Mobile       string    `json:"mobile" gorm:"column:mobile"`
}

func NewUser() *User {
	return new(User)
}
func (*User) TableName() string {
	return "user"
}
func (a *User) Insert(ctx context.Context) {
	err := core.DB.Create(&a).Error
	if err != nil {
		logger.Info(ctx, "db insert fail", err)
	}
}
func (a *User) Save(ctx context.Context) {
	core.DB.Save(a)
}

//获取分页列表
func (a *User) PageList(where map[string]interface{}, page int, pagesize int, order string) (result []User) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pagesize

	core.DB.Where(where).Order(order).Offset(offset).Limit(pagesize).Find(&result)
	return
}
func (a *User) Del(ctx context.Context) {
	core.DB.Delete(a)
}

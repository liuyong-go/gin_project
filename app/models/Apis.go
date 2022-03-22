package models

import (
	"context"
	"time"

	"github.com/liuyong-go/gin_project/libs/logger"
)

type Apis struct {
	ID                         int       `json:"id"`
	SiteId                     int       `json:"site_id"`
	DepartmentId               string    `json:"department_id"`
	Title                      string    `json:"title"`
	RequestPath                string    `json:"request_path"`
	LeaderId                   int       `json:"leader_id"`
	ManagerId                  int       `json:"manager_id"`
	SyapiId                    int       `json:"syapi_id"`
	Comment                    string    `json:"comment"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	Level                      int       `json:"level"`
	AvgExecTimeOfYesterday     float32   `json:"avg_exec_time_of_yesterday"`
	SuccessPercentOfYesterday  float32   `json:"success_percent_of_yesterday"`
	RequestCountYesterday      int       `json:"request_count_yesterday"`
	AvgExecTimeByAlert         float32   `json:"avg_exec_time_by_alert"`
	SuccessPercentByAlert      float32   `json:"success_percent_by_alert"`
	AvgSlowExecTimeOfYesterday float32   `json:"avg_slow_exec_time_of_yesterday"`
}

func NewApis() *Apis {
	return &Apis{}
}

func (*Apis) TableName() string {
	return "apis"
}
func (a *Apis) Create(ctx context.Context) {
	err := DB.Create(&a).Error
	if err != nil {
		logger.Info(ctx, "db insert fail", err)
	}
}
func (a *Apis) Get(ctx context.Context) Apis {
	var api Apis
	DB.First(&api)
	return api
}

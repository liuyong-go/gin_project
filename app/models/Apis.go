package models

import "time"

type Apis struct {
	ID                         string    `json:"id"`
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

func (*Apis) TableName() string {
	return "apis"
}

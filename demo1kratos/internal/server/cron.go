package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/orzkratos/cronkratos"
	"github.com/orzkratos/demokratos/demo1kratos/internal/service"
	"github.com/robfig/cron/v3"
)

// NewCronServer creates a new cron server and registers cron jobs
// 创建新的 cron server 并注册定时任务
func NewCronServer(cronSvc *service.CronService, logger log.Logger) *cronkratos.Server {
	srv := cronkratos.NewServer(cron.New(), logger)
	cronkratos.RegisterCronServer(srv, cronSvc)
	return srv
}

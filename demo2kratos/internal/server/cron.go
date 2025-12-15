package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/orzkratos/cronkratos"
	"github.com/orzkratos/demokratos/demo2kratos/internal/service"
	"github.com/robfig/cron/v3"
)

// NewCronServer creates a new cron server and registers cron jobs with locker
// 创建新的 cron server 并注册带锁的定时任务
func NewCronServer(cronSvc *service.CronService, logger log.Logger) *cronkratos.Server {
	srv := cronkratos.NewServer(cron.New(), logger)
	cronkratos.RegisterCronServerL(srv, cronSvc)
	return srv
}

package service

import (
	"context"

	"github.com/orzkratos/demokratos/demo1kratos/internal/biz"
	"github.com/robfig/cron/v3"
)

// CronService handles cron job registration
// 处理定时任务注册
type CronService struct {
	task *biz.TaskUsecase
}

// NewCronService creates a new CronService instance
// 创建新的 CronService 实例
func NewCronService(task *biz.TaskUsecase) *CronService {
	return &CronService{task: task}
}

// RegisterCron registers cron jobs to cron instance
// 注册定时任务到 cron 实例
func (s *CronService) RegisterCron(ctx context.Context, c *cron.Cron) {
	// Sync data every minute
	// 每分钟同步数据
	_, _ = c.AddFunc("* * * * *", func() {
		_ = s.task.SyncData(ctx)
	})

	// Cleanup data every hour
	// 每小时清理数据
	_, _ = c.AddFunc("0 * * * *", func() {
		_ = s.task.CleanupData(ctx)
	})
}

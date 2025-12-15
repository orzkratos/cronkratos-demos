package service

import (
	"context"
	"sync"

	"github.com/orzkratos/demokratos/demo2kratos/internal/biz"
	"github.com/robfig/cron/v3"
)

// CronService handles cron job registration with locker
// 处理带锁的定时任务注册
type CronService struct {
	task *biz.TaskUsecase
}

// NewCronService creates a new CronService instance
// 创建新的 CronService 实例
func NewCronService(task *biz.TaskUsecase) *CronService {
	return &CronService{task: task}
}

// RegisterCron registers cron jobs with locker for graceful shutdown
// 注册定时任务，使用锁支持优雅退出
func (s *CronService) RegisterCron(ctx context.Context, c *cron.Cron, locker sync.Locker) {
	// Sync data every minute with lock protection
	// 每分钟同步数据，带锁保护
	_, _ = c.AddFunc("* * * * *", func() {
		locker.Lock()
		defer locker.Unlock()
		_ = s.task.SyncData(ctx)
	})

	// Cleanup data every hour with lock protection
	// 每小时清理数据，带锁保护
	_, _ = c.AddFunc("0 * * * *", func() {
		locker.Lock()
		defer locker.Unlock()
		_ = s.task.CleanupData(ctx)
	})
}

package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// TaskUsecase handles scheduled task business logic
// 处理定时任务的业务逻辑
type TaskUsecase struct {
	log *log.Helper
}

// NewTaskUsecase creates a new TaskUsecase instance
// 创建新的 TaskUsecase 实例
func NewTaskUsecase(logger log.Logger) *TaskUsecase {
	return &TaskUsecase{
		log: log.NewHelper(logger),
	}
}

// SyncData performs data synchronization task
// 执行数据同步任务
func (uc *TaskUsecase) SyncData(ctx context.Context) error {
	uc.log.WithContext(ctx).Infof("SyncData task executed at %s", time.Now().Format(time.RFC3339))
	// TODO: Add actual sync logic here
	// TODO: 在这里添加实际的同步逻辑
	return nil
}

// CleanupData performs data cleanup task
// 执行数据清理任务
func (uc *TaskUsecase) CleanupData(ctx context.Context) error {
	uc.log.WithContext(ctx).Infof("CleanupData task executed at %s", time.Now().Format(time.RFC3339))
	// TODO: Add actual cleanup logic here
	// TODO: 在这里添加实际的清理逻辑
	return nil
}

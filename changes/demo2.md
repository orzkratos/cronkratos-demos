# Changes

Code differences compared to source project demokratos.

## cmd/demo2kratos/main.go (+3 -1)

```diff
@@ -11,6 +11,7 @@
 	"github.com/go-kratos/kratos/v2/middleware/tracing"
 	"github.com/go-kratos/kratos/v2/transport/grpc"
 	"github.com/go-kratos/kratos/v2/transport/http"
+	"github.com/orzkratos/cronkratos"
 	"github.com/orzkratos/demokratos/demo2kratos/internal/conf"
 	"github.com/yyle88/done"
 	"github.com/yyle88/must"
@@ -31,7 +32,7 @@
 	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
 }
 
-func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
+func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, cs *cronkratos.Server) *kratos.App {
 	return kratos.New(
 		kratos.ID(done.VCE(os.Hostname()).Omit()),
 		kratos.Name(Name),
@@ -41,6 +42,7 @@
 		kratos.Server(
 			gs,
 			hs,
+			cs,
 		),
 	)
 }
```

## cmd/demo2kratos/wire_gen.go (+5 -1)

```diff
@@ -2,6 +2,7 @@
 
 //go:generate go run -mod=mod github.com/google/wire/cmd/wire
 //go:build !wireinject
+// +build !wireinject
 
 package main
 
@@ -28,7 +29,10 @@
 	greeterService := service.NewGreeterService(greeterUsecase)
 	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
 	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
-	app := newApp(logger, grpcServer, httpServer)
+	taskUsecase := biz.NewTaskUsecase(logger)
+	cronService := service.NewCronService(taskUsecase)
+	cronkratosServer := server.NewCronServer(cronService, logger)
+	app := newApp(logger, grpcServer, httpServer, cronkratosServer)
 	return app, func() {
 		cleanup()
 	}, nil
```

## internal/biz/biz.go (+1 -1)

```diff
@@ -3,4 +3,4 @@
 import "github.com/google/wire"
 
 // ProviderSet is biz providers.
-var ProviderSet = wire.NewSet(NewGreeterUsecase)
+var ProviderSet = wire.NewSet(NewGreeterUsecase, NewTaskUsecase)
```

## internal/biz/task.go (+36 -0)

```diff
@@ -0,0 +1,36 @@
+package biz
+
+import (
+	"context"
+	"time"
+
+	"github.com/go-kratos/kratos/v2/log"
+)
+
+// TaskUsecase handles scheduled task business logic
+// 处理定时任务的业务逻辑
+type TaskUsecase struct {
+	log *log.Helper
+}
+
+// NewTaskUsecase creates a new TaskUsecase instance
+// 创建新的 TaskUsecase 实例
+func NewTaskUsecase(logger log.Logger) *TaskUsecase {
+	return &TaskUsecase{
+		log: log.NewHelper(logger),
+	}
+}
+
+// SyncData performs data synchronization task
+// 执行数据同步任务
+func (uc *TaskUsecase) SyncData(ctx context.Context) error {
+	uc.log.WithContext(ctx).Infof("SyncData task executed at %s", time.Now().Format(time.RFC3339))
+	return nil
+}
+
+// CleanupData performs data cleanup task
+// 执行数据清理任务
+func (uc *TaskUsecase) CleanupData(ctx context.Context) error {
+	uc.log.WithContext(ctx).Infof("CleanupData task executed at %s", time.Now().Format(time.RFC3339))
+	return nil
+}
```

## internal/server/cron.go (+16 -0)

```diff
@@ -0,0 +1,16 @@
+package server
+
+import (
+	"github.com/go-kratos/kratos/v2/log"
+	"github.com/orzkratos/cronkratos"
+	"github.com/orzkratos/demokratos/demo2kratos/internal/service"
+	"github.com/robfig/cron/v3"
+)
+
+// NewCronServer creates a new cron server and registers cron jobs with locker
+// 创建新的 cron server 并注册带锁的定时任务
+func NewCronServer(cronSvc *service.CronService, logger log.Logger) *cronkratos.Server {
+	srv := cronkratos.NewServer(cron.New(), logger)
+	cronkratos.RegisterCronServerL(srv, cronSvc)
+	return srv
+}
```

## internal/server/server.go (+1 -1)

```diff
@@ -5,4 +5,4 @@
 )
 
 // ProviderSet is server providers.
-var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)
+var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewCronServer)
```

## internal/service/cron.go (+41 -0)

```diff
@@ -0,0 +1,41 @@
+package service
+
+import (
+	"context"
+	"sync"
+
+	"github.com/orzkratos/demokratos/demo2kratos/internal/biz"
+	"github.com/robfig/cron/v3"
+)
+
+// CronService handles cron job registration with locker
+// 处理带锁的定时任务注册
+type CronService struct {
+	task *biz.TaskUsecase
+}
+
+// NewCronService creates a new CronService instance
+// 创建新的 CronService 实例
+func NewCronService(task *biz.TaskUsecase) *CronService {
+	return &CronService{task: task}
+}
+
+// RegisterCron registers cron jobs with locker for graceful shutdown
+// 注册定时任务，使用锁支持优雅退出
+func (s *CronService) RegisterCron(ctx context.Context, c *cron.Cron, locker sync.Locker) {
+	// Sync data every minute with lock protection
+	// 每分钟同步数据，带锁保护
+	_, _ = c.AddFunc("* * * * *", func() {
+		locker.Lock()
+		defer locker.Unlock()
+		_ = s.task.SyncData(ctx)
+	})
+
+	// Cleanup data every hour with lock protection
+	// 每小时清理数据，带锁保护
+	_, _ = c.AddFunc("0 * * * *", func() {
+		locker.Lock()
+		defer locker.Unlock()
+		_ = s.task.CleanupData(ctx)
+	})
+}
```

## internal/service/service.go (+1 -1)

```diff
@@ -3,4 +3,4 @@
 import "github.com/google/wire"
 
 // ProviderSet is service providers.
-var ProviderSet = wire.NewSet(NewGreeterService)
+var ProviderSet = wire.NewSet(NewGreeterService, NewCronService)
```


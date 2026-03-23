package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ops-timer-backend/internal/api/handler"
	"ops-timer-backend/internal/api/router"
	"ops-timer-backend/internal/config"
	"ops-timer-backend/internal/model"
	"ops-timer-backend/internal/pkg/auth"
	"ops-timer-backend/internal/pkg/email"
	pkgoauth "ops-timer-backend/internal/pkg/oauth"
	"ops-timer-backend/internal/pkg/scheduler"
	"ops-timer-backend/internal/repository"
	"ops-timer-backend/internal/service"

	"go.uber.org/zap"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	initAdmin := flag.Bool("init-admin", false, "初始化/重置管理员账户")
	adminUser := flag.String("admin-user", "admin", "管理员用户名")
	adminPass := flag.String("admin-pass", "admin123", "管理员密码")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	zapLogger := initLogger(cfg.Log)
	defer zapLogger.Sync()

	db := initDatabase(cfg.Database)

	autoMigrate(db)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	unitRepo := repository.NewUnitRepository(db)
	unitLogRepo := repository.NewUnitLogRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	todoGroupRepo := repository.NewTodoGroupRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	// Auth
	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpiryHours)

	// Email
	emailSvc := email.NewService(&cfg.SMTP)
	if emailSvc.Enabled() {
		zapLogger.Info("SMTP 邮件通知已启用", zap.String("host", cfg.SMTP.Host))
	} else {
		zapLogger.Info("SMTP 邮件通知未配置，跳过邮件功能")
	}

	// OAuth / OIDC
	var oauthSvc *pkgoauth.Service
	if cfg.OAuth.IsConfigured() {
		var oauthErr error
		oauthSvc, oauthErr = pkgoauth.NewService(&cfg.OAuth)
		if oauthErr != nil {
			zapLogger.Warn("OAuth 初始化失败，OAuth 登录不可用", zap.Error(oauthErr))
		} else {
			zapLogger.Info("OAuth OIDC 已启用", zap.String("issuer", cfg.OAuth.IssuerURL))
		}
	} else {
		zapLogger.Info("OAuth 未配置，跳过 OAuth 登录功能")
	}

	// Services
	authService := service.NewAuthService(userRepo, jwtManager, &cfg.Auth)
	unitService := service.NewUnitService(unitRepo, unitLogRepo)
	projectService := service.NewProjectService(projectRepo, unitRepo)
	todoService := service.NewTodoService(todoRepo, todoGroupRepo)
	notifService := service.NewNotificationService(notifRepo)

	if err := authService.EnsureAdminExists(*adminUser, *adminPass); err != nil {
		zapLogger.Fatal("创建管理员账户失败", zap.Error(err))
	}

	if *initAdmin {
		zapLogger.Info("管理员账户已初始化", zap.String("username", *adminUser))
		os.Exit(0)
	}

	// Handlers
	authHandler := handler.NewAuthHandler(authService, emailSvc)
	oauthHandler := handler.NewOAuthHandler(oauthSvc, authService)
	unitHandler := handler.NewUnitHandler(unitService)
	projectHandler := handler.NewProjectHandler(projectService, unitService)
	todoHandler := handler.NewTodoHandler(todoService)
	notifHandler := handler.NewNotificationHandler(notifService)

	// Router
	r := router.NewRouter(&router.RouterConfig{
		AuthHandler:    authHandler,
		OAuthHandler:   oauthHandler,
		UnitHandler:    unitHandler,
		ProjectHandler: projectHandler,
		TodoHandler:    todoHandler,
		NotifHandler:   notifHandler,
		JWTManager:     jwtManager,
		AuthService:    authService,
		Logger:         zapLogger,
		CorsOrigins:    cfg.Server.CorsOrigins,
	})

	engine := r.Setup()

	// Scheduler
	sched := scheduler.NewScheduler(unitRepo, notifRepo, userRepo, emailSvc, zapLogger)
	if err := sched.Start(cfg.Scheduler.NotificationScanInterval); err != nil {
		zapLogger.Error("启动定时任务失败", zap.Error(err))
	}
	defer sched.Stop()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	zapLogger.Info("计时器服务已启动", zap.String("addr", addr))

	if err := engine.Run(addr); err != nil {
		zapLogger.Fatal("服务启动失败", zap.Error(err))
	}
}

func initLogger(cfg config.LogConfig) *zap.Logger {
	var zapCfg zap.Config
	if cfg.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	switch cfg.Level {
	case "debug":
		zapCfg.Level.SetLevel(zap.DebugLevel)
	case "warn":
		zapCfg.Level.SetLevel(zap.WarnLevel)
	case "error":
		zapCfg.Level.SetLevel(zap.ErrorLevel)
	default:
		zapCfg.Level.SetLevel(zap.InfoLevel)
	}

	l, err := zapCfg.Build()
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	return l
}

func initDatabase(cfg config.DatabaseConfig) *gorm.DB {
	dir := filepath.Dir(cfg.DSN)
	if dir != "." {
		os.MkdirAll(dir, 0755)
	}

	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	return db
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Unit{},
		&model.UnitLog{},
		&model.TodoGroup{},
		&model.Todo{},
		&model.Notification{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
}

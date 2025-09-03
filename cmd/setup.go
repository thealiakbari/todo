package cmd

import (
	"context"

	todoItemHttpAdaptor "github.com/thealiakbari/todoapp/internal/adapters/inbound/http/todo"
	todoItemOutboundRepo "github.com/thealiakbari/todoapp/internal/adapters/outbound/db/pg"
	todoItemApp "github.com/thealiakbari/todoapp/internal/application/todo"
	todoItemService "github.com/thealiakbari/todoapp/internal/domain/todo"
	todoInterface "github.com/thealiakbari/todoapp/internal/ports/inbound/todo"
	todoItemRepo "github.com/thealiakbari/todoapp/internal/ports/outbound/todo"
	"github.com/thealiakbari/todoapp/pkg/common/config"
	"github.com/thealiakbari/todoapp/pkg/common/db"
	"github.com/thealiakbari/todoapp/pkg/common/i18next"
	"github.com/thealiakbari/todoapp/pkg/common/logger"
	"golang.org/x/text/language"
)

type RepositoryStorage struct {
	todoItemRepo todoItemRepo.TodoItemRepository
}

type ServiceStorage struct {
	todoItemSvc todoInterface.TodoItemService
}

type ApplicationStorage struct {
	todoItemApp todoItemApp.TodoItemHttpApp
}

type HttpAdaptorStorage struct {
	TodoItemAdaptor todoItemHttpAdaptor.Adaptor
}

type SetupConfig struct {
	Ctx                context.Context
	Conf               *config.AppConfig
	Logger             logger.Logger
	DB                 db.DBWrapper
	HttpAdaptorStorage HttpAdaptorStorage
}

func Setup() *SetupConfig {
	ctx := context.Background()
	conf := config.LoadConfig("./config/todoapp.yml")

	log, err := logger.New(
		conf.Mode,
		conf.ServiceName,
		"todoapp",
	)
	if err != nil {
		panic(err)
	}

	err = i18next.NewLanguage(language.Make(conf.Language))
	if err != nil {
		panic(err)
	}

	logInfra := log.CloneAsInfra()
	err = db.Migrate(conf.DB.Postgres, logInfra)
	if err != nil {
		logInfra.Panicf("Migration failed: %s\n", err.Error())
	}
	logInfra.Info("Migrations successfully done.")

	gormDB, err := db.NewPostgresConn(ctx, conf.DB.Postgres)
	if err != nil {
		panic(err)
	}

	dbw := db.NewDBWrapper(gormDB)

	repos := NewRepositoryStorage(dbw)
	services := NewServiceStorage(log, repos)

	httpApps := NewHttpAppStorage(dbw, services)
	httpAdaptors := NewHttpAdaptorStorage(httpApps)

	return &SetupConfig{
		Ctx:                ctx,
		Conf:               conf,
		Logger:             log,
		DB:                 dbw,
		HttpAdaptorStorage: httpAdaptors,
	}
}

func NewHttpAppStorage(
	db db.DBWrapper,
	services ServiceStorage,
) ApplicationStorage {
	return ApplicationStorage{
		todoItemApp: todoItemApp.NewTodoItemHttpApp(services.todoItemSvc, db),
	}
}

func NewRepositoryStorage(db db.DBWrapper) RepositoryStorage {
	return RepositoryStorage{
		todoItemRepo: todoItemOutboundRepo.NewTodoItemRepository(db),
	}
}

func NewServiceStorage(log logger.Logger, repos RepositoryStorage) ServiceStorage {
	return ServiceStorage{
		todoItemSvc: todoItemService.NewTodoItemService(todoItemService.TodoItemConfig{Logger: log, TodoItemRepo: repos.todoItemRepo}),
	}
}

func NewHttpAdaptorStorage(
	httpApps ApplicationStorage,
) HttpAdaptorStorage {
	return HttpAdaptorStorage{
		TodoItemAdaptor: todoItemHttpAdaptor.Adaptor{TodoItemHttpApp: httpApps.todoItemApp},
	}
}

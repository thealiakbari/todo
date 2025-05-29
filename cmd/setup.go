package cmd

import (
	"context"
	coreEnum "github.com/thealiakbari/hichapp/pkg/core/enum"

	pollApp "github.com/thealiakbari/hichapp/app/poll/service"
	tagApp "github.com/thealiakbari/hichapp/app/tag/service"
	userApp "github.com/thealiakbari/hichapp/app/user/service"
	pollService "github.com/thealiakbari/hichapp/internal/poll"
	pollRepo "github.com/thealiakbari/hichapp/internal/poll/domain/repository"
	tagService "github.com/thealiakbari/hichapp/internal/tag"
	tagRepo "github.com/thealiakbari/hichapp/internal/tag/domain/repository"
	userService "github.com/thealiakbari/hichapp/internal/user"
	userRepo "github.com/thealiakbari/hichapp/internal/user/domain/repository"
	"github.com/thealiakbari/hichapp/pkg/common/config"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/i18next"
	"github.com/thealiakbari/hichapp/pkg/common/kafka"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	"github.com/thealiakbari/hichapp/pkg/common/store"
	pollHttpAdaptor "github.com/thealiakbari/hichapp/pkg/infrastructure/ports/in_bounds/http/poll"
	tagHttpAdaptor "github.com/thealiakbari/hichapp/pkg/infrastructure/ports/in_bounds/http/tag"
	userHttpAdaptor "github.com/thealiakbari/hichapp/pkg/infrastructure/ports/in_bounds/http/user"
	"golang.org/x/text/language"
)

type RepositoryStorage struct {
	pollRepo pollRepo.PollRepository
	tagRepo  tagRepo.TagRepository
	userRepo userRepo.UserRepository
}

type ServiceStorage struct {
	userSvc userService.User
	pollSvc pollService.Poll
	tagSvc  tagService.Tag
}

type BusinessFlowStorage struct{}

type HttpAppStorage struct {
	pollApp pollApp.PollHttpApp
	tagApp  tagApp.TagHttpApp
	userApp userApp.UserHttpApp
}

type HttpAdaptorStorage struct {
	UserAdaptor userHttpAdaptor.Adaptor
	PollAdaptor pollHttpAdaptor.Adaptor
	TagAdaptor  tagHttpAdaptor.Adaptor
}

type OutboundStorage struct{}

type ConsumerStorage struct {
	PollConsumers map[coreEnum.Topic]kafka.HandlerFn
}

type SetupConfig struct {
	Ctx                context.Context
	Conf               *config.AppConfig
	Logger             logger.Logger
	DB                 db.DBWrapper
	Kafka              kafka.Kafka
	Store              store.Store
	HttpAdaptorStorage HttpAdaptorStorage
	ConsumerStorage    ConsumerStorage
}

func Setup() *SetupConfig {
	ctx := context.Background()
	conf := config.LoadConfig("./config/hichapp.yml")

	log, err := logger.New(
		conf.Mode,
		conf.ServiceName,
		"hichapp",
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

	storeSvc := store.New(
		conf.DB.Redis.Address,
		conf.DB.Redis.Password,
		conf.DB.Redis.DB,
		logInfra,
	)
	err = storeSvc.Ping(ctx)
	if err != nil {
		panic(err)
	}

	dbw := db.NewDBWrapper(gormDB)

	kaf := kafka.NewKafka(conf.Kafka, log.CloneAsInfra(), dbw)

	repos := NewRepositoryStorage(dbw)
	services := NewServiceStorage(log, repos, dbw, conf)
	//authInfraSvc := authInfra.NewAuthInfraHTTP(
	//	log,
	//	conf.Core.Auth,
	//	storeSvc,
	//	services.userSvc,
	//)

	flows := NewBusinessFlowStorage(conf, dbw, log, storeSvc, services)
	httpApps := NewHttpAppStorage(dbw, services, flows, log, conf)
	httpAdaptors := NewHttpAdaptorStorage(log, dbw, httpApps)
	consumers := NewConsumerStorage(log, dbw, services)

	return &SetupConfig{
		Ctx:                ctx,
		Conf:               conf,
		Logger:             log,
		DB:                 dbw,
		Kafka:              kaf,
		Store:              storeSvc,
		HttpAdaptorStorage: httpAdaptors,
		ConsumerStorage:    consumers,
	}
}

func NewHttpAppStorage(
	db db.DBWrapper,
	services ServiceStorage,
	flows BusinessFlowStorage,
	log logger.Logger,
	cfg *config.AppConfig,
) HttpAppStorage {
	return HttpAppStorage{
		userApp: userApp.NewUserHttpApp(services.userSvc, db),
		pollApp: pollApp.NewPollHttpApp(services.pollSvc, db),
		tagApp:  tagApp.NewTagHttpApp(services.tagSvc, db),
	}
}

func NewRepositoryStorage(db db.DBWrapper) RepositoryStorage {
	return RepositoryStorage{
		userRepo: userRepo.NewUserRepository(db),
		pollRepo: pollRepo.NewPollRepository(db),
		tagRepo:  tagRepo.NewTagRepository(db),
	}
}

func NewServiceStorage(log logger.Logger, repos RepositoryStorage, db db.DBWrapper, cfg *config.AppConfig) ServiceStorage {
	return ServiceStorage{
		userSvc: userService.NewUserService(userService.UserConfig{Logger: log, UserRepo: repos.userRepo}),
		pollSvc: pollService.NewPollService(pollService.PollConfig{Logger: log, PollRepo: repos.pollRepo}),
		tagSvc:  tagService.NewTagService(tagService.TagConfig{Logger: log, TagRepo: repos.tagRepo}),
	}
}

func NewBusinessFlowStorage(
	cfg *config.AppConfig,
	dbw db.DBWrapper,
	log logger.Logger,
	storeSvc store.Store,
	services ServiceStorage,
) BusinessFlowStorage {
	return BusinessFlowStorage{}
}

func NewHttpAdaptorStorage(
	log logger.Logger,
	db db.DBWrapper,
	httpApps HttpAppStorage,
) HttpAdaptorStorage {
	return HttpAdaptorStorage{
		PollAdaptor: pollHttpAdaptor.Adaptor{PollHttpApp: httpApps.pollApp},
		TagAdaptor:  tagHttpAdaptor.Adaptor{TagHttpApp: httpApps.tagApp},
		UserAdaptor: userHttpAdaptor.Adaptor{UserHttpApp: httpApps.userApp},
	}
}

func NewConsumerStorage(
	log logger.Logger,
	db db.DBWrapper,
	services ServiceStorage,
) ConsumerStorage {
	return ConsumerStorage{
		//PollConsumers: pollApp.(services.userSvc, log, db),
	}
}

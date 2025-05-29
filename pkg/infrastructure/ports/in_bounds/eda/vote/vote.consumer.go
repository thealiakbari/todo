package cross

import (
	userService "github.com/thealiakbari/hichapp/internal/poll"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/kafka"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	coreEnum "github.com/thealiakbari/hichapp/pkg/core/enum"
)

func NewAdaptor(userSvc userService.Poll, log logger.Logger, db db.DBWrapper) map[coreEnum.Topic]kafka.HandlerFn {
	consumer := service.NewConsumer(userSvc, log, db)
	handlers := make(map[coreEnum.Topic]kafka.HandlerFn, 0)
	handlers[enum.UmsUserCrudEvent] = consumer.OnUmsUserCrudEvent

	return handlers
}

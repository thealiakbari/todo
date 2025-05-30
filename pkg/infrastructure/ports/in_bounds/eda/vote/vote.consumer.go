package cross

import (
	"github.com/thealiakbari/hichapp/app/poll/service"
	"github.com/thealiakbari/hichapp/app/user/domain/enum"
	pollService "github.com/thealiakbari/hichapp/internal/poll"
	"github.com/thealiakbari/hichapp/pkg/common/db"
	"github.com/thealiakbari/hichapp/pkg/common/kafka"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
	coreEnum "github.com/thealiakbari/hichapp/pkg/core/enum"
)

func NewAdaptor(pollSvc pollService.Poll, log logger.Logger, db db.DBWrapper) map[coreEnum.Topic]kafka.HandlerFn {
	consumer := service.NewConsumer(pollSvc, log, db)
	handlers := make(map[coreEnum.Topic]kafka.HandlerFn, 0)
	handlers[enum.HichAppVoteActionEvent] = consumer.OnHichAppVoteEvent

	return handlers
}

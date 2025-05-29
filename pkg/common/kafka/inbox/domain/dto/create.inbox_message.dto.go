package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/inbox/domain/enum"
)

type InboxMessageCreateDto struct {
	Id            string           `validate:"uuid4"`
	AggregateId   string           `validate:"uuid4"`
	Type          string           `validate:"required"`
	CorrelationId string           `validate:"uuid4"`
	TraceId       string           `validate:"uuid4"`
	State         enum.InboxState  `validate:"required"`
	Status        enum.InboxStatus `validate:"required"`
	RetryCount    int              `validate:"required,gte=0"`
	WaitDuration  *int
	Metadata      *string
	Payload       *string
}

func (i InboxMessageCreateDto) Validate() error {
	validate := validator.New()
	err := validate.Struct(i)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}

	return nil
}

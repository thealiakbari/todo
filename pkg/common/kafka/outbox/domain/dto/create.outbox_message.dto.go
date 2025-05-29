package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/thealiakbari/hichapp/pkg/common/kafka/outbox/domain/enum"
)

type OutboxMessageCreateDto struct {
	TraceId       string                 `validate:"required,uuid4"`
	AggregateId   string                 `validate:"required,uuid4"`
	AggregateType string                 `validate:"required"`
	Type          enum.OutboxMessageType `validate:"required"`
	Name          string                 `validate:"required"`
	Payload       *string                `validate:"required"`
}

func (i OutboxMessageCreateDto) Validate() error {
	validate := validator.New()
	err := validate.Struct(i)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}

	return nil
}

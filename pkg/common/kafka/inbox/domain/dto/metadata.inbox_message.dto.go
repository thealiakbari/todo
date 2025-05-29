package dto

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type InboxMessageHeaderDto struct {
	OutboxId      string `validate:"required,uuid4"`
	TraceId       string `validate:"required,uuid4"`
	CorrelationId string `validate:"required,uuid4"`
	Type          string `validate:"required"`
	Name          string `validate:"required"`
	Timestamp     int64  `validate:"required"`
}

func (i InboxMessageHeaderDto) Validate() error {
	validate := validator.New()
	err := validate.Struct(i)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}

	return nil
}

func LoadInboxMessageHeader(headers map[string]string) (*InboxMessageHeaderDto, error) {
	ihs := InboxMessageHeaderDto{
		OutboxId:      headers["id"],
		TraceId:       headers["tid"],
		CorrelationId: headers["cid"],
		Type:          headers["type"],
		Name:          headers["name"],
	}

	ts, err := strconv.ParseInt(headers["ts"], 10, 64)
	if err != nil {
		return nil, errors.New("field \"ts\" :" + err.Error())
	}

	ihs.Timestamp = ts

	err = ihs.Validate()
	if err != nil {
		return nil, err
	}

	return &ihs, nil
}

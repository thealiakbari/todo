package kafka

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type (
	HandlerFn = func(context.Context, message.Payload) error
)

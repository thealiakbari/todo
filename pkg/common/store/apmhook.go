package store

import (
	"bytes"
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.elastic.co/apm/v2"
)

const ApmStatKeyPrefix = "redis_key_"

type apmStat struct {
	Opt string `json:"operation"`
	Key string `json:"key"`
}

type hook struct{}

func newApmHook() redis.Hook {
	return &hook{}
}

func (r *hook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (r *hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		span, _ := apm.StartSpanOptions(ctx, getCmdName(cmd), "db.redis", apm.SpanOptions{
			ExitSpan: true,
		})
		defer span.End()
		return next(ctx, cmd)
	}
}

func (r *hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		var cmdNameBuf bytes.Buffer
		for i, cmd := range cmds {
			if i != 0 {
				cmdNameBuf.WriteString(", ")
			}
			cmdNameBuf.WriteString(getCmdName(cmd))
		}
		span, _ := apm.StartSpanOptions(ctx, cmdNameBuf.String(), "db.redis", apm.SpanOptions{
			ExitSpan: true,
		})
		defer span.End()
		return next(ctx, cmds)
	}
}

func getCmdName(cmd redis.Cmder) string {
	cmdName := strings.ToUpper(cmd.Name())
	if cmdName == "" {
		cmdName = "(empty command)"
	}
	return cmdName
}

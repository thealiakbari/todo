package kafka

import (
	"fmt"
	"sort"
	"strings"

	wm "github.com/ThreeDotsLabs/watermill"
	"github.com/thealiakbari/hichapp/pkg/common/logger"
)

type watermillLogger struct {
	logger logger.InfraLogger
	fields wm.LogFields
}

func newWatermillLogger(logger logger.InfraLogger) watermillLogger {
	return watermillLogger{
		logger: logger,
	}
}

func (l watermillLogger) fieldsStringer(fields wm.LogFields) string {
	fieldsStr := ""

	allFields := l.fields.Add(fields)

	keys := make([]string, len(allFields))
	i := 0
	for field := range allFields {
		keys[i] = field
		i++
	}

	sort.Strings(keys)

	for _, key := range keys {
		var valueStr string
		value := allFields[key]

		if stringer, ok := value.(fmt.Stringer); ok {
			valueStr = stringer.String()
		} else {
			valueStr = fmt.Sprintf("%v", value)
		}

		if strings.Contains(valueStr, " ") {
			valueStr = `"` + valueStr + `"`
		}

		fieldsStr += key + "=" + valueStr + " "
	}

	return fieldsStr
}

func (l watermillLogger) Error(msg string, err error, fields wm.LogFields) {
	if l.logger == nil {
		return
	}
	l.logger.Errorf(msg+"%s, err: %s", l.fieldsStringer(fields), err.Error())
}

func (l watermillLogger) Info(msg string, fields wm.LogFields) {
	if l.logger == nil {
		return
	}
	l.logger.Infof(msg+" %s", l.fieldsStringer(fields))
}

func (l watermillLogger) Debug(msg string, fields wm.LogFields) {
	if l.logger == nil {
		return
	}
	l.logger.Infof(msg+" %s", l.fieldsStringer(fields))
}

func (l watermillLogger) Trace(msg string, fields wm.LogFields) {
	if l.logger == nil {
		return
	}
	l.logger.Infof(msg+" %s", l.fieldsStringer(fields))
}

func (l watermillLogger) With(fields wm.LogFields) wm.LoggerAdapter {
	return watermillLogger{
		logger: l.logger,
		fields: fields,
	}
}

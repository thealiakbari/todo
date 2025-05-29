package logger

import (
	"log/slog"
	"time"

	"go.uber.org/zap"
)

type FieldType = uint

const (
	UnknownType FieldType = iota
	BinaryType
	StringType
	IntType
	TimeType
	ErrorType

	modeDev = "development"
)

const (
	serviceTag        = "service"
	prefixTagZapField = "prefixTagZapField"
	zapFieldPrefix    = "zapfield_prefix"
	errorFieldName    = "_error"
	errorKeyName      = "error"
)

type Field struct {
	Key   string
	Type  FieldType
	Value interface{}
}

func (f *Field) toZapField() zap.Field {
	switch f.Type {
	case BinaryType:
		return zap.Binary(f.Key, f.Value.([]byte))
	case StringType:
		return zap.String(f.Key, f.Value.(string))
	case IntType:
		return zap.Int(f.Key, f.Value.(int))
	case TimeType:
		return zap.Time(f.Key, f.Value.(time.Time))
	case ErrorType:
		return zap.Error(f.Value.(error))
	default:
		return zap.Any(f.Key, f.Value)
	}
}

func (f *Field) toSlogAttr() slog.Attr {
	switch f.Type {
	case BinaryType:
		return slog.Any(f.Key, f.Value.([]byte))
	case StringType:
		return slog.String(f.Key, f.Value.(string))
	case IntType:
		return slog.Int(f.Key, f.Value.(int))
	case TimeType:
		return slog.Time(f.Key, f.Value.(time.Time))
	case ErrorType:
		return slog.Any(f.Key, f.Value.(error))
	default:
		return slog.Any(f.Key, f.Value)
	}
}

func String(key string, value string) Field {
	return Field{
		Key:   key,
		Type:  StringType,
		Value: value,
	}
}

func Error(err error) Field {
	return Field{
		Key:   errorKeyName,
		Type:  ErrorType,
		Value: err,
	}
}

func Int(key string, value int) Field {
	return Field{
		Key:   key,
		Type:  IntType,
		Value: value,
	}
}

func Time(key string, value time.Time) Field {
	return Field{
		Key:   key,
		Type:  TimeType,
		Value: value,
	}
}

func Binary(key string, value []byte) Field {
	return Field{
		Key:   key,
		Type:  BinaryType,
		Value: value,
	}
}

func Any(key string, value interface{}) Field {
	return Field{
		Key:   key,
		Type:  UnknownType,
		Value: value,
	}
}

func fieldsToInterface(fields []Field) []interface{} {
	res := make([]interface{}, 0, len(fields))
	for _, f := range fields {
		res = append(res, f.toZapField())
	}
	return res
}

func fieldsToSlogAttrs(fields []Field) []slog.Attr {
	res := make([]slog.Attr, 0, len(fields))
	for _, f := range fields {
		res = append(res, f.toSlogAttr())
	}
	return res
}

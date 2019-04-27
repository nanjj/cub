package logs

import (
	"fmt"
	"math"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap/zapcore"
)

type SpanZapCore struct {
	*FieldsSpan
	zapcore.LevelEnabler
}

type Context interface {
	Fields() []log.Field
}

type FieldsSpan struct {
	Context
	span   opentracing.Span
	fields []zapcore.Field
}

func (c *FieldsSpan) Fields() (fields []log.Field) {
	fields = append(fields, logFields(c.fields)...)
	if c.Context == nil {
		return
	}
	fields = append(fields, c.Context.Fields()...)
	return
}

func NewFieldsSpan(sp opentracing.Span) *FieldsSpan {
	return &FieldsSpan{span: sp}
}

func NewSpanZapCore(sp opentracing.Span) *SpanZapCore {
	return &SpanZapCore{NewFieldsSpan(sp), zapcore.DebugLevel}
}

func (c *SpanZapCore) With(fields []zapcore.Field) (core zapcore.Core) {
	if len(fields) == 0 {
		return c
	}
	parent := c.FieldsSpan
	span := parent.span
	core = &SpanZapCore{
		FieldsSpan:   &FieldsSpan{parent, span, fields},
		LevelEnabler: c.LevelEnabler,
	}
	return core
}

func (c *SpanZapCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *SpanZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) (err error) {
	span := c.span
	if level := entry.Level; level >= zapcore.ErrorLevel {
		span.SetTag("error", true)
	}
	spanFields := entryFields(entry)
	if len(fields) != 0 {
		spanFields = append(spanFields, logFields(fields)...)
	}

	if withFields := c.Fields(); len(withFields) != 0 {
		spanFields = append(spanFields, withFields...)
	}
	span.LogFields(spanFields...)
	return
}

func (c *SpanZapCore) Sync() error {
	c.span.Finish()
	return nil
}

func entryFields(entry zapcore.Entry) (fields []log.Field) {
	fields = make([]log.Field, 3)
	fields[0] = log.String("message", entry.Message)
	fields[1] = log.String("caller", entry.Caller.TrimmedPath())
	if stack := entry.Stack; stack != "" {
		fields[3] = log.String("stack", stack)
	} else {
		fields = fields[0:2]
	}
	return
}

func logFields(zfields []zapcore.Field) (fields []log.Field) {
	fields = make([]log.Field, len(zfields))
	for i := range zfields {
		fields[i] = logField(zfields[i])
	}
	return
}

func logField(zapField zapcore.Field) log.Field {
	switch zapField.Type {

	case zapcore.BoolType:
		val := false
		if zapField.Integer >= 1 {
			val = true
		}
		return log.Bool(zapField.Key, val)
	case zapcore.Float32Type:
		return log.Float32(zapField.Key, math.Float32frombits(uint32(zapField.Integer)))
	case zapcore.Float64Type:
		return log.Float64(zapField.Key, math.Float64frombits(uint64(zapField.Integer)))
	case zapcore.Int64Type:
		return log.Int64(zapField.Key, int64(zapField.Integer))
	case zapcore.Int32Type:
		return log.Int32(zapField.Key, int32(zapField.Integer))
	case zapcore.StringType:
		return log.String(zapField.Key, zapField.String)
	case zapcore.StringerType:
		return log.String(zapField.Key, zapField.Interface.(fmt.Stringer).String())
	case zapcore.Uint64Type:
		return log.Uint64(zapField.Key, uint64(zapField.Integer))
	case zapcore.Uint32Type:
		return log.Uint32(zapField.Key, uint32(zapField.Integer))
	case zapcore.DurationType:
		return log.String(zapField.Key, time.Duration(zapField.Integer).String())
	case zapcore.ErrorType:
		return log.Error(zapField.Interface.(error))
	default:
		return log.Object(zapField.Key, zapField.Interface)
	}
}

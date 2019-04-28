package logs_test

//go:generate mockgen -destination mock_$GOFILE -package logs_test github.com/opentracing/opentracing-go Span,Tracer,SpanContext

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/nanjj/cub/logs"
	"go.uber.org/zap"

	gomock "github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go/log"
)

type Matches func(v interface{}) bool

func (f Matches) Matches(x interface{}) bool {
	return f(x)
}

func (f Matches) String() string {
	return "Matches Dismatch"
}

func TestNewSpanZapCore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	span := NewMockSpan(ctrl)
	_, _, line, _ := runtime.Caller(0)
	gomock.InOrder(
		span.EXPECT().LogFields(
			log.String("message", "something wrong"),
			log.String("caller", fmt.Sprintf("logs/zap_test.go:%d", line+17)),
			log.String("AccountID", "45e0d4be68f711e991a32fc786358b81")).
			Times(1),
		span.EXPECT().LogFields(
			log.String("message", "ID Created"),
			log.String("caller", fmt.Sprintf("logs/zap_test.go:%d", line+21)),
			log.String("AccountID", "45e0d4be68f711e991a32fc786358b81"),
			log.String("ID", "e655cd3c-68f7-11e9-90c0-174f343089a5")).
			Times(1))
	span.EXPECT().Finish().Times(1)
	span.EXPECT().SetTag("error", true).Times(1)
	logger := zap.New(logs.NewSpanZapCore(span),
		zap.AddCaller())
	logger.Error("something wrong",
		zap.String("AccountID", "45e0d4be68f711e991a32fc786358b81"))
	logger.
		With(zap.String("ID", "e655cd3c-68f7-11e9-90c0-174f343089a5")).
		Info("ID Created",
			zap.String("AccountID", "45e0d4be68f711e991a32fc786358b81"))
	defer logger.Sync()
}

func TestDebugLevel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	span := NewMockSpan(ctrl)
	span.EXPECT().Finish().Times(1)
	logger := zap.New(logs.NewSpanZapCore(span))
	logger.Debug("hello")
	defer logger.Sync()
}

func TestCallStack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	span := NewMockSpan(ctrl)
	span.EXPECT().Finish().Times(1)
	span.EXPECT().LogFields(log.String("message", "enter"),
		gomock.AssignableToTypeOf(log.Field{})).Times(1)
	logger := zap.New(logs.NewSpanZapCore(span), zap.AddStacktrace(zap.InfoLevel))
	defer logger.Sync()
	logger.Info("enter")
}

func TestWithFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	span := NewMockSpan(ctrl)
	span.EXPECT().Finish().Times(1)
	gomock.InOrder(
		span.EXPECT().LogFields(log.String("message", "Step 1"),
			log.Bool("enter", true)).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 2"),
			log.Float32("weight", 13.2)).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 3"),
			log.Float64("height", 82.1)).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 4"),
			log.Int64("count", 100)).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 5"),
			log.Int32("count", 100)).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 6"),
			log.Uint64("count", 100)).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 7"),
			log.Uint32("count", 100)).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 8"),
			log.String("time", time.Second.String())).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 9"),
			log.Error(fmt.Errorf("something wrong"))).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 10"),
			log.String("stringer", time.Second.String())).Times(1),
		span.EXPECT().LogFields(log.String("message", "Step 11"),
			log.Object("complex", (12+0i))).Times(1),
	)
	logger := zap.New(logs.NewSpanZapCore(span))
	defer logger.Sync()
	logger.With().With(zap.Bool("enter", true)).Info("Step 1")
	logger.With(zap.Float32("weight", 13.2)).Info("Step 2")
	logger.With(zap.Float64("height", 82.1)).Info("Step 3")
	logger.With(zap.Int64("count", 100)).Info("Step 4")
	logger.With(zap.Int32("count", 100)).Info("Step 5")
	logger.With(zap.Uint64("count", 100)).Info("Step 6")
	logger.With(zap.Uint32("count", 100)).Info("Step 7")
	logger.With(zap.Duration("time", time.Second)).Info("Step 8")
	logger.With(zap.Error(fmt.Errorf("something wrong"))).Info("Step 9")
	logger.With(zap.Stringer("stringer", time.Second)).Info("Step 10")
	logger.With(zap.Complex128("complex", 12+0i)).Info("Step 11")
}

func TestSetLogsLevel(t *testing.T) {
	logs.Level = -2
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	span := NewMockSpan(ctrl)
	span.EXPECT().Finish().Times(1)
	span.EXPECT().LogFields(log.String("message", "what's wrong?")).Times(1)
	os.Setenv("LOGS_LEVEL", "debug")
	defer os.Unsetenv("LOGS_LEVEL")
	logger := zap.New(logs.NewSpanZapCore(span))
	defer logger.Sync()
	logger.Debug("what's wrong?")
}

func TestNewSpanLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	span := NewMockSpan(ctrl)
	span.EXPECT().Finish().Times(1)
	span.EXPECT().LogFields(log.String("message", "print stack trace"),
		Matches(func(x interface{}) bool {
			if v, ok := x.(log.Field); ok {
				if v.Key() == "stack" {
					return true
				}
			}
			return false
		})).Times(1)
	logger := logs.NewSpanLogger(span)
	logger = logger.With()
	logger = logger.With(zap.Stack("stack"))
	logger.Info("print stack trace")
	defer logger.Finish()
}

func TestWithOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	span := NewMockSpan(ctrl)
	span.EXPECT().Finish().Times(1)
	span.EXPECT().LogFields(log.String("message", "print caller"),
		Matches(func(x interface{}) bool {
			if field, ok := x.(log.Field); ok {
				if field.Key() == "caller" {
					return true
				}
			}
			return false
		})).Times(1)
	logger := logs.NewSpanLogger(span)
	logger = logger.WithOptions()
	logger = logger.WithOptions(zap.AddCaller())
	defer logger.Finish()
	logger.Info("print caller")
}

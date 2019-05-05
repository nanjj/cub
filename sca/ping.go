package sca

import (
	"context"
	"time"

	"github.com/nanjj/cub/logs"
	"go.uber.org/zap"
)

// Ping
func (r *Runner) Ping(ctx context.Context, req Payload) (rep Payload, err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Ping")
	defer sp.Finish()
	sp.Info("Ping", zap.String("name", r.Name()))
	t := DataObject{}
	err = t.Encode(time.Now().UTC())
	if err == nil {
		rep = Payload{DataObject(r.Name()), t}
	}
	return
}

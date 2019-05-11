package sca

import (
	"context"
	"time"

	"github.com/nanjj/cub/logs"
	"github.com/nanjj/cub/sdo"
	"go.uber.org/zap"
)

// Ping
func (r *Runner) Ping(ctx context.Context, req sdo.Payload) (rep sdo.Payload, err error) {
	sp, ctx := logs.StartSpanFromContext(ctx, "Ping")
	defer sp.Finish()
	sp.Info("Ping", zap.String("name", r.Name()))
	t := sdo.DataObject{}
	err = t.Encode(time.Now().UTC())
	if err == nil {
		rep = sdo.Payload{sdo.DataObject(r.Name()), t}
	}
	return
}

package sca

import (
	"context"

	"nanomsg.org/go/mangos/v2"
)

func Send(ctx context.Context, sock mangos.Socket, e *Event) (err error) {
	sp, ctx := StartSpanFromContext(ctx, "Send")
	defer sp.Finish()
	err = e.Emit(ctx, sock)
	return
}

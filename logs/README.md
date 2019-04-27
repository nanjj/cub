# Logs
```
 __        ______     _______      _______.
|  |      /  __  \   /  _____|    /       |
|  |     |  |  |  | |  |  __     |   (----`
|  |     |  |  |  | |  | |_ |     \   \
|  `----.|  `--'  | |  |__| | .----)   |
|_______| \______/   \______| |_______/

`logs` send logs to opentracing span via zap logger.

```
## Install

```
go get github.com/nanjj/cub/logs
```

## Import
```
import github.com/nanjj/cub/logs

```

## Usage

```
logger, ctx := logs.StartSpanFromContext(ctx, "Create")
defer logger.Finish() // Use opentracing.Span API to finish span

logger.Info("Enter") // Use zap logger api to send log
logger.SetTag("RequestID", requestID) // Use opentracing.Span API to set tag
```
## Notice

1. Do not send too many logs in a span, for Jaeger users, the UDP
   packet limit may be met for too many logs in a span.
2. Defer call `logger.Finish()` or `logger.Sync()`, no need to do both.

package metric

import "context"

type Recorder interface {
	Count(ctx context.Context, name string, value int64)
}

package antsqueue

import "context"

type TJob struct {
	Key  string
	Func func(ctx context.Context)
	Ctx  context.Context
}

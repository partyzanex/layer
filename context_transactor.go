package layer

import "context"

type ContextTransactor interface {
	Begin(ctx context.Context) (context.Context, error)
	Execute(ctx context.Context, err error) error
}

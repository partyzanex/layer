package layer

import (
	"context"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
)

type BoilExecutor interface {
	boil.ContextBeginner
	boil.ContextExecutor
}

type contextKey struct {
	KeyName string
}

var (
	ExecutorKey = contextKey{
		KeyName: "executor",
	}
	TransactorKey = contextKey{
		KeyName: "transactor",
	}
)

func GetExecutor(ctx context.Context, executor BoilExecutor) (context.Context, boil.ContextExecutor) {
	if ctx != nil {
		executor, exists := ctx.Value(ExecutorKey).(boil.ContextExecutor)
		if exists {
			return ctx, executor
		}
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, ExecutorKey, executor), executor
}

func GetTransactor(ctx context.Context) (context.Context, boil.ContextTransactor) {
	if ctx != nil {
		transactor, exists := ctx.Value(TransactorKey).(boil.ContextTransactor)
		if exists {
			return ctx, transactor
		}
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	return ctx, nil
}

func ExecuteTransaction(transactor boil.ContextTransactor, err error) error {
	var errTr error

	if err == nil {
		errTr = transactor.Commit()
	}

	if errTr != nil {
		errTr = transactor.Rollback()
		if errTr != nil {
			errTr = errors.Wrap(errTr, "rollback failed")
		}
	}

	return errTr
}

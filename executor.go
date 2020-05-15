package layer

import (
	"context"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

	ErrCreateTransaction  = errors.New("creating transaction failed")
	ErrContextIsNil       = errors.New("context is nil")
	ErrTransactorIsExists = errors.New("transactor is exists")
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

func CreateTransaction(ctx context.Context, beginner boil.ContextBeginner) (context.Context, boil.ContextTransactor, error) {
	if ctx == nil {
		return context.TODO(), nil, ErrContextIsNil
	}

	transactor, exists := ctx.Value(TransactorKey).(boil.ContextTransactor)
	if exists {
		return ctx, transactor, ErrTransactorIsExists
	}

	tx, err := beginner.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "creating transaction failed")
	}

	ctx = context.WithValue(ctx, ExecutorKey, tx)
	ctx = context.WithValue(ctx, TransactorKey, tx)

	return ctx, tx, nil
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

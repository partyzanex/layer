package layer_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"

	"github.com/partyzanex/layer"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
)

const (
	driverName = "postgres"
	dsn        = "dbname=lp host=localhost port=5432 user=postgres password=535353 sslmode=disable"
)

func TestGetExecutor(t *testing.T) {
	t.Parallel()
	ex, err := getTestExecutor()
	if err != nil {
		t.Fatalf("getting test executor failed: %s", err)
	}

	ctx, e := layer.GetExecutor(nil, nil)
	if e != nil {
		t.Errorf("getting executor failed: expected nil, got %v", e)
	}
	if ctx == context.Background() {
		t.Errorf("getting context failed: expected %v, got %v", context.Background(), ctx)
	}

	ctx = context.WithValue(ctx, layer.ExecutorKey, ex)

	_, e = layer.GetExecutor(ctx, ex)
	if _, ok := e.(layer.BoilExecutor); !ok && e != ex {
		t.Errorf("wrong executor: expected %v, got %v", ex, e)
	}
}

func TestGetTransactor(t *testing.T) {
	t.Parallel()
	ex, err := getTestExecutor()
	if err != nil {
		t.Fatalf("getting test executor failed: %s", err)
	}

	ctx, tr := layer.GetTransactor(nil)
	if tr != nil {
		t.Errorf("getting transactor failed: expected nil, got %v", tr)
	}
	if ctx != context.TODO() {
		t.Errorf("getting context failed: expected %v, got %v", context.Background(), ctx)
	}

	tr, err = ex.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("creating transaction failed: %s", err)
	}
	if _, ok := tr.(boil.ContextTransactor); !ok {
		t.Fatalf("transaction is not implements boil.ContextTransactor %v", tr)
	}
	if _, ok := tr.(boil.ContextExecutor); !ok {
		t.Fatalf("transaction is not implements boil.ContextExecutor %v", tr)
	}
	defer tr.Rollback()

	ctx = context.WithValue(ctx, layer.TransactorKey, tr)

	_, tx := layer.GetTransactor(ctx)
	if tx == nil {
		t.Errorf("transactor is nil")
	}
	if _, ok := tx.(boil.ContextExecutor); !ok {
		t.Errorf("transaction is not implements boil.ContextExecutor %v", tx)
	}
	if _, ok := tx.(boil.ContextTransactor); !ok {
		t.Errorf("transaction is not implements boil.ContextExecutor %v", tr)
	}
	if tx != tr {
		t.Errorf("getting transactor failed: expected %v, got %v", tr, tx)
	}
}

func TestExecuteTransaction(t *testing.T) {
	t.Parallel()
	var err error

	ex, err := getTestExecutor()
	if err != nil {
		t.Fatalf("getting test executor failed: %s", err)
	}

	ctx := context.TODO()

	tr, err := ex.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("creating transaction failed")
	}

	result, err := tr.Exec("select now() as dt")
	if err != nil {
		t.Errorf("executing query failed: %s", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		t.Errorf("getting affected rows failed: %s", err)
	}
	if n != 1 {
		t.Errorf("wrong count of affected rows: expected 1, got %d", n)
	}

	err = layer.ExecuteTransaction(tr, err)
	if err != nil {
		t.Errorf("executing transaction failed: %s", err)
	}
}

func getTestExecutor() (layer.BoilExecutor, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, errors.Wrap(err, "create connection failed")
	}

	return db, nil
}

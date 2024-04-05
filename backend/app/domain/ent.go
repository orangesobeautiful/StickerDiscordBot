package domain

import (
	"context"

	"backend/app/ent"
	"backend/app/pkg/hserr"

	"golang.org/x/xerrors"
)

type BaseEntRepoInterface interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

var _ BaseEntRepoInterface = (*BaseEntRepo)(nil)

type BaseEntRepo struct {
	entCli *ent.Client
}

func NewBaseEntRepo(client *ent.Client) *BaseEntRepo {
	return &BaseEntRepo{entCli: client}
}

func (r *BaseEntRepo) WithTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx := ctxGetTx(ctx)
	if tx != nil {
		return fn(ctx)
	}

	tx, err := r.getNewClient().Tx(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "create tx")
	}
	ctx = ctxSetTx(ctx, tx)
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(ctx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = xerrors.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return hserr.NewInternalError(err, "committing transaction")
	}
	return nil
}

type txKeyType struct{}

func (r *BaseEntRepo) GetEntClient(ctx context.Context) *ent.Client {
	tx := ctxGetTx(ctx)
	if tx != nil {
		return tx.Client()
	}

	return r.getNewClient()
}

func (r *BaseEntRepo) getNewClient() *ent.Client {
	return r.entCli.Debug()
}

func ctxGetTx(ctx context.Context) *ent.Tx {
	txAny := ctx.Value(txKeyType{})
	if txAny == nil {
		return nil
	}

	return txAny.(*ent.Tx)
}

func ctxSetTx(ctx context.Context, tx *ent.Tx) context.Context {
	return context.WithValue(ctx, txKeyType{}, tx)
}

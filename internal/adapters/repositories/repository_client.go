package repositories

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

const (
	DBTransactionCtx = "_app_db_trx_ctx"
)

type DBClient interface {
	Begin(ctx context.Context) context.Context
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error

	client(ctx context.Context) *gorm.DB
}

type dbClient struct {
	gormDB *gorm.DB
}

func NewDBClient(gormDB *gorm.DB) DBClient {
	return &dbClient{
		gormDB: gormDB,
	}
}

func (db *dbClient) Begin(ctx context.Context) context.Context {
	if transactionCtx := ctx.Value(DBTransactionCtx); transactionCtx != nil {
		if _, ok := transactionCtx.(*gorm.DB); ok {
			return ctx
		}
	}

	dbTrxNewSession := db.gormDB.Begin()
	return context.WithValue(ctx, DBTransactionCtx, dbTrxNewSession)
}

func (db *dbClient) Rollback(ctx context.Context) error {
	if transactionCtx := ctx.Value(DBTransactionCtx); transactionCtx != nil {
		if dbTrx, ok := transactionCtx.(*gorm.DB); ok {
			return dbTrx.Rollback().Error
		}
	}

	return fmt.Errorf("rollback fail | db transaction context is invalid")
}

func (db *dbClient) Commit(ctx context.Context) error {
	if transactionCtx := ctx.Value(DBTransactionCtx); transactionCtx != nil {
		if dbTrx, ok := transactionCtx.(*gorm.DB); ok {
			return dbTrx.Commit().Error
		}
	}

	return fmt.Errorf("commit fail | db transaction context is invalid")
}

func (db *dbClient) client(ctx context.Context) *gorm.DB {
	if transactionCtx := ctx.Value(DBTransactionCtx); transactionCtx != nil {
		if dbTrx, ok := transactionCtx.(*gorm.DB); ok {
			return dbTrx.WithContext(ctx)
		}
	}

	return db.gormDB.WithContext(ctx)
}

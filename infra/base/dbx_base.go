package base

import (
	"context"
	"database/sql"
	"github.com/tietang/dbx"
	"log"
)

const TX = "tx"

type BaseDao struct {
	TX *sql.Tx
}

func (d *BaseDao) SetTx(tx *sql.Tx) {
	d.TX = tx
}

type txFunc func(runner *dbx.TxRunner) error

// Tx 事务执行帮助函数，简化代码
func Tx(fn func(*dbx.TxRunner) error) error {
	return TxContext(context.Background(), fn)
}

// TxContext 事务执行帮助函数，简化代码，需要传入上下文
func TxContext(ctx context.Context, fn func(runner *dbx.TxRunner) error) error {
	return DbxDatabase().Tx(fn)
}

// WithValueContext 将runner绑定到上下文，并创建一个新的WithValueContext
func WithValueContext(parent context.Context, runner *dbx.TxRunner) context.Context {
	return context.WithValue(parent, TX, runner)
}

func ExecuteContext(ctx context.Context, fn func(*dbx.TxRunner) error) error {
	tx, ok := ctx.Value(TX).(*dbx.TxRunner)
	if !ok || tx == nil {
		log.Panic("是否在事务函数块中使用？")
	}
	return fn(tx)
}

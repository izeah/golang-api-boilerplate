package trxmanager

import (
	"fmt"
	"runtime/debug"

	"boilerplate/internal/abstraction"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type trxManager struct {
	db *gorm.DB
}

type trxFn func(ctx *abstraction.Context) error

func New(db *gorm.DB) *trxManager {
	return &trxManager{db}
}

func (g *trxManager) WithTrx(pCtx *abstraction.Context, fn trxFn) (err error) {
	if pCtx.Trx != nil {
		return fn(pCtx)
	}

	pCtx.Trx = &abstraction.TrxContext{
		Db: g.db.WithContext(pCtx.Request().Context()).Clauses(dbresolver.Write).Begin(),
	}

	fmt.Printf("\nBEGIN\n")

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			pCtx.Trx.Db.Rollback()
			fmt.Printf("\nROLLBACK\n\n")
			logrus.Error(p)
			err = fmt.Errorf("panic happened because: %s, stacktrace: %s", fmt.Sprintf("%v", p), string(debug.Stack()))
		} else if err != nil {
			// error occurred, rollback
			pCtx.Trx.Db.Rollback()
			fmt.Printf("\nROLLBACK\n\n")
		} else {
			// all good, commit
			err = pCtx.Trx.Db.Commit().Error
			fmt.Printf("\nCOMMIT\n\n")
		}
		pCtx.Trx = nil
	}()

	return fn(pCtx)
}

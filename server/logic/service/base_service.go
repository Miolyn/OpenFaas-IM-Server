package service

import (
	"OpenFaaS-Logic/dal/db"
	"github.com/gin-gonic/gin"
)

type Base struct {
	Affair *db.AffairDB
}

func NewBase() *Base {
	return &Base{}
}

type AffairInterface interface {
	AffairImInit() bool
	AffairAppInit() bool
	AffairBegin(c *gin.Context) func()
	AffairRollbackIfError(err error) bool
	AffairFinished() bool
	AffairInitWithAffair(a *db.AffairDB)
	AffairRollback()
}

func (b *Base) AffairImInit() bool {
	b.Affair = new(db.AffairDB)
	err := b.Affair.NewImAffairs()
	if err != nil {
		return true
	}
	return false
}

func (b *Base) AffairInitWithAffair(a *db.AffairDB) {
	b.Affair.InitAffairsWithTx(a.GetTX())
}

func (b *Base) AffairBegin(c *gin.Context) func() {
	return b.Affair.DeferFunc(c)
}

func (b *Base) AffairRollbackIfError(err error) bool {
	if err != nil {
		b.Affair.RollBackIfError(err)
		return true
	}
	return false
}

func (b *Base) AffairRollback() {
	b.Affair.RollBack()
}

func (b *Base) AffairFinished() bool {
	err := b.Affair.Commit()
	if err != nil {
		return true
	}
	return false
}

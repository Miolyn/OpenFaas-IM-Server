package db

import (
	"OpenFaaS-Logic/pkg/e"
	"OpenFaaS-Logic/pkg/logger"
	"OpenFaaS-Logic/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AffairDB struct {
	tx     *gorm.DB
	isInit bool
}

func (a *AffairDB) NewImAffairs() (err error) {
	db := DBConn.FindersIMDB
	a.tx = db.Begin()
	if err = a.tx.Error; err != nil {
		return
	}
	a.isInit = true
	return nil
}

func (a *AffairDB) InitAffairsWithTx(tx *gorm.DB) {
	a.tx = tx
	a.isInit = true
}
func (a *AffairDB) DeferFunc(c *gin.Context) func() {
	if a.tx == nil {
		return func() {
			panic("affair not init")
		}
	}
	return func() {
		if r := recover(); r != nil {
			a.tx.Rollback()
			logger.Logger.Error("mysql affair error", r)
			response.ErrorHandler(c, e.MYSQL_AFFAIR_ERROR)
		}
	}
}

func (a *AffairDB) RollBackIfError(err error) {
	if err != nil {
		a.tx.Rollback()
	}
}
func (a *AffairDB) RollBack() {
	a.tx.Rollback()
}
func (a *AffairDB) Commit() (err error) {
	return a.tx.Commit().Error
}

func (a *AffairDB) GetTX() *gorm.DB {
	return a.tx
}

func (a *AffairDB) IsInit() bool {
	return a.isInit
}

package db

import (
	"OpenFaas-User/pkg/logger"
	"OpenFaas-User/st"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type dbConn struct {
	FindersIMDB *gorm.DB
}

var (
	DBConn   *dbConn
	dbConfig DataBase
)

type DataBase struct {
	DriverName   string `yaml:"DriverName"`
	DBHost       string `yaml:"DBHost"`
	DBPort       string `yaml:"DBPort"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
	Timeout      string `yaml:"Timeout"`      //Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	ReadTimeout  string `yaml:"ReadTimeout"`  //Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	WriteTimeout string `yaml:"WriteTimeout"` //Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	User         string `yaml:"User"`
	Password     string `yaml:"Password"`
	DBName       string `yaml:"DBName"`
	Config       string `yaml:"Config"`
}

func init() {
	dbConfig := DataBase{
		DriverName:   "mysql2",
		DBHost:       "49.4.114.179",
		DBPort:       "3306",
		MaxIdleConns: 100,
		MaxOpenConns: 500,
		Timeout:      "400ms",
		ReadTimeout:  "200ms",
		WriteTimeout: "400ms",
		User:         "root",
		Password:     "123456",
		DBName:       "imdb",
		Config:       "charset=utf8mb4&parseTime=True&loc=Local",
	}
	DBConn = &dbConn{
		FindersIMDB: getConn(dbConfig),
	}
}
func getConn(database DataBase) *gorm.DB {
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?%s", database.User, database.Password, database.DBHost, database.DBPort, database.DBName, database.Config)
	if db, err := gorm.Open("mysql", connStr); err != nil {
		logger.Logger.Fatalf("%s 链接异常 :%v", database.DBName, err)
		return nil
	} else {
		st.Debug("db init success")
		return db
	}
}

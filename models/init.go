package models

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jinzhu/gorm"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
	"github.com/wangyi1310/mycloud-disk/pkg/util"
)

const (
	SQLITE  = "sqlite"
	SQLITE3 = "sqlite3"
	MYSQL   = "mysql"
	MSSQL   = "mssql"
	UNSET   = "unset"
)

var DB *gorm.DB

var DB_CONNECT_FUNC_MAP = map[string]func(string) (*gorm.DB, error){
	SQLITE:  connectSQLite,
	MYSQL:   connectMySQL,
	MSSQL:   connectMySQL,
	SQLITE3: connectSQLite,
	UNSET:   connectSQLite,
}

func connectSQLite(confDBType string) (*gorm.DB, error) {
	dbFilePath := ""
	if conf.DatabaseConfig.DBFile == "" {
		dbFilePath = "sqlite.db"
	} else {
		dbFilePath = conf.DatabaseConfig.DBFile
	}
	dbFullFilePath := util.RelativePath(dbFilePath)
	return gorm.Open("sqlite", util.RelativePath(dbFullFilePath))
}

func connectMySQL(confDBType string) (*gorm.DB, error) {
	var host string
	if conf.DatabaseConfig.UnixSocket {
		host = fmt.Sprintf("unix(%s)",
			conf.DatabaseConfig.Host)
	} else {
		host = fmt.Sprintf("(%s:%d)",
			conf.DatabaseConfig.Host,
			conf.DatabaseConfig.Port)
	}

	db, err := gorm.Open(confDBType, fmt.Sprintf("%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local",
		conf.DatabaseConfig.User,
		conf.DatabaseConfig.Password,
		host,
		conf.DatabaseConfig.Name,
		conf.DatabaseConfig.Charset))
	return db, err
}

func Init() {
	log.Log().Info("Initializing database connection...")
	var (
		db     *gorm.DB
		err    error
		dbType = strings.ToLower(conf.DatabaseConfig.Type)
	)
	connectFunc, ok := DB_CONNECT_FUNC_MAP[dbType]
	if !ok {
		log.Log().Panic("Unsupport database type %q.", dbType)
	}

	db, err = connectFunc(dbType)
	if err != nil {
		log.Log().Panic("Failed to connect database: %s", err)
	}

	if conf.SystemConfig.Debug {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	db.DB().SetMaxIdleConns(50)
	if dbType == SQLITE || dbType == SQLITE3 || dbType == UNSET {
		db.DB().SetMaxOpenConns(1)
	} else {
		db.DB().SetMaxOpenConns(100)
	}

	db.DB().SetConnMaxLifetime(time.Second * 30)
	DB = db

	migration()
}

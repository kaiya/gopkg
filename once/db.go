package once

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenDB(dbName string, migrateModel ...interface{}) (db *gorm.DB) {
	//db
	sqlUser := os.Getenv("SQL_USER")
	sqlPass := os.Getenv("SQL_PASS")
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000)/%s?parseTime=true", sqlUser, sqlPass, dbName)))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(migrateModel...)
	return
}

package startup

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/mysqljob/repository/dao"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(43.154.97.245:13316)/webook_job"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

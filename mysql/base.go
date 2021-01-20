package mysql

import (
	"fmt"
	"fund/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	mysqlConfig := config.Load("mysql")
	Db, _ = gorm.Open("mysql", fmt.Sprintf(
		"%v:%v@%v(%v)/%v?charset=utf8&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["network"],
		mysqlConfig["address"]+":"+mysqlConfig["port"],
		mysqlConfig["database"]))
}

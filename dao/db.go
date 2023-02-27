package dao

import (
	"fmt"
	//引入 sqlx
	"github.com/jmoiron/sqlx"
	"log"
	//项目内的 配置信息
	"awesomeProject/conf"
)

var BaseDb sqlx.DB

// app.go 里面调用
func InitDB() {

	datasource := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8", conf.Config.Mysql.User, conf.Config.Mysql.Passwd, conf.Config.Mysql.Host, conf.Config.Mysql.DBName)
	db, err := sqlx.Open("mysql", datasource)
	if err != nil {
		log.Fatal(err)
		return
	}
	BaseDb = *db

}

// app.go 里面调用
func CloseDB() {
	err := BaseDb.Close()
	if err != nil {
		log.Fatal(err)
	}
}

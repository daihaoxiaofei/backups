package core

import (
	"backups/config"
	"database/sql"
	"fmt"
	"testing"
)

func Test_(t *testing.T) {
	for _, TM := range config.C.Tables {
		go tes(TM)
	}
	select {}
}

func tes(TM config.Table) {
	fmt.Println(TM)
}

func Test_checkLink(t *testing.T) {
	checkLink()
}

func Test_sqlLink(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True",
		`root`, `jdys123.`, `39.108.36.31`, `xmh`)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("数据库打开错误")
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.Ping()
	if err != nil {
		panic("数据库连接错误")
	}
}

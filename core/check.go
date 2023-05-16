package core

import (
	"backups/config"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
)

// 检查文件夹创建情况
func checkDir() {
	for _, TM := range config.C.Tables {
		fileDir := path + TM.TName
		s, err := os.Stat(fileDir)
		if err != nil || !s.IsDir() {
			err = os.MkdirAll(fileDir, os.ModePerm)
			if err != nil {
				log.Println(`路径创建失败:` + err.Error())
				return
			} else {
				log.Println(`路径创建: ` + fileDir)
			}
		}
		log.Println(`备份数据库: `, TM.Host, TM.TName)
	}
}

// 检查数据库连通性
func checkLink() {
	wg := sync.WaitGroup{}
	wg.Add(len(config.C.Tables))
	for _, table := range config.C.Tables {
		go func(table config.Table) {
			defer wg.Done()
			dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True",
				table.UName, table.Passwd, table.Host, table.TName)

			db, err := sql.Open("mysql", dsn)
			if err != nil {
				log.Printf("测试数据库打开错误 %s -> %s %s", table.TName, table.Host, err.Error())
				return
			}
			defer db.Close()
			err = db.Ping()
			if err != nil {
				log.Printf("测试数据库连接错误 %s -> %s %s", table.TName, table.Host, err.Error())
				return
			}
			log.Printf("测试数据库连接成功 %s -> %s", table.TName, table.Host)
		}(table)
	}
	wg.Wait()
}

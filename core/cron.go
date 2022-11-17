package core

import (
	"backups/config"
	"github.com/robfig/cron/v3"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	format = `2006-01-02 15_04_05` // 显示格式
	path = `./backups/`            // 备份目录
)

func Cron() {
	checkDir()
	c := cron.New()

	if config.C.Cron == `` {
		// 每小时一次 给个随机的分钟数
		rand.Seed(time.Now().UnixNano()) // 随机数
		Minute := rand.Intn(59)          // 返回[0,n)
		config.C.Cron = strconv.Itoa(Minute) + " * * * *"
	}
	_, err := c.AddFunc(config.C.Cron, func() {
		backup()
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(`开始备份程序`)
	c.Start()
	select {}
}

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

func backup() {
	// 备份
	for _, TM := range config.C.Tables {
		go Dump(TM)
	}
	// 清理备份
	for _, TM := range config.C.Tables {
		go ClearFile(TM.TName)
	}
}

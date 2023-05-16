package core

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"

	"backups/config"
)

var (
	format = `2006-01-02 15_04_05` // 显示格式
	path   = `./backups/`          // 备份目录
)

func Cron() {
	checkDir()
	checkLink()
	c := cron.New()

	if config.C.Cron == `` {
		// 每小时一次 给个随机的分钟数
		rand.Seed(time.Now().UnixNano()) // 随机数
		Minute := rand.Intn(59)          // 返回[0,n)
		config.C.Cron = strconv.Itoa(Minute) + " * * * *"
	}
	_, err := c.AddFunc(config.C.Cron, func() {
		// 备份
		for _, TM := range config.C.Tables {
			go Dump(TM)
		}
		// 清理备份
		for _, TM := range config.C.Tables {
			go ClearFile(TM.TName)
		}
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(`开始备份程序`)
	c.Start()
	select {}
}

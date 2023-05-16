package core

import (
	"backups/config"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var saveHour int

func init() {
	saveHour = config.C.Hour
	if saveHour < 0 || saveHour > 24 {
		saveHour = 6
	}
}

// 返回可留日期
func saveDate() []string {
	t := time.Now()
	// 当前小时
	onHour := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	// 当日最佳时刻
	onDay := time.Date(t.Year(), t.Month(), t.Day(), saveHour, 0, 0, 0, t.Location())
	// 当月1号最佳时刻
	onMonth := time.Date(t.Year(), t.Month(), 1, saveHour, 0, 0, 0, t.Location())

	return []string{
		onHour.Format(format), // 当前小时 当天保留
		onHour.Add(-time.Hour).Format(format),
		onHour.Add(-time.Hour * 2).Format(format),
		onHour.Add(-time.Hour * 4).Format(format),
		onHour.Add(-time.Hour * 8).Format(format),
		onHour.Add(-time.Hour * 16).Format(format),

		onDay.Format(format), // 当日凌晨6点 近三天保留
		onDay.AddDate(0, 0, -1).Format(format),
		onDay.AddDate(0, 0, -2).Format(format),

		onMonth.Format(format), // 当月1号凌晨6点 近三月保留
		onMonth.AddDate(0, -1, 0).Format(format),
		onMonth.AddDate(0, -2, 0).Format(format),
		onMonth.AddDate(0, -3, 0).Format(format),
	}
}

// ClearFile 清理文件 关键在于筛选和删除
func ClearFile(TName string) {
	filePath := path + TName
	// 扫描文件
	var files, err = ioutil.ReadDir(filePath)
	if err != nil {
		log.Println(`ReadDir错误: `, err)
		return
	}

	// 筛选 返回可留日期
	saveDateList := saveDate()
	for _, file := range files {
		// 以.sql结尾 && 文件是否要删除
		if strings.HasSuffix(file.Name(), `.sql`) && isDelete(file.Name(), saveDateList) {
			err = os.Remove(filePath + `/` + file.Name()) // 删除
			if err != nil {
				log.Println(`删除文件出错: `, err)
			} else {
				log.Println(`删除文件: `, filePath+`/`+file.Name())
			}
		}
	}
}

// 判断文件是否要删除
func isDelete(fileName string, saveDateList []string) bool {
	for _, dateName := range saveDateList {
		// 只判断前13位 即到小时级别
		if fileName[0:13] == dateName[0:13] {
			return false
		}
	}
	return true
}

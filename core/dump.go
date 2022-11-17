package core

import (
	"backups/config"
	"log"
	"os"
	"os/exec"
	"time"
)

// 备份
func Dump(TM config.Table) {
	// 文件
	date := time.Now().Format(format)
	file, err := os.OpenFile(path+TM.TName+"/"+date+".sql",
		os.O_CREATE|os.O_RDWR,
		os.ModePerm|os.ModeTemporary)
	log.Println(`备份文件: `, path+TM.TName+"/"+date+".sql")
	if err != nil {
		log.Println(`打开sql文件失败: `, err)
		return
	}
	defer file.Close()

	argv := []string{
		"-h" + TM.Host,
		"-u" + TM.UName,
		"-p" + TM.Passwd,
		TM.TName,
	}
	cmd := exec.Command("mysqldump", argv...)

	cmd.Stdout = file
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		log.Println(`错误Start`, err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		log.Println(`错误Wait`, err)
		return
	}
}

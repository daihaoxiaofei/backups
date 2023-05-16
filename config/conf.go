package config

import (
	"backups/aes"
	"crypto/tls"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Table struct {
	Host   string `yaml:"Host"`
	TName  string `yaml:"TName"`
	UName  string `yaml:"UName"`
	Passwd string `yaml:"Passwd"`
}

type c struct {
	Url    string  `yaml:"Url"`
	CKey   string  `yaml:"CKey"`
	Cron   string  `yaml:"Cron"` // 备份时间规则
	Hour   int     `yaml:"Hour"` // 保留的最佳时刻 如凌晨6点
	Tables []Table `yaml:"Tables"`
}

var C c

func init() {
	confPath := "config.yaml"

	notMain := false
	switch runtime.GOOS {
	case `linux`:
		if os.Args[0][len(os.Args[0])-5:] == `.test` {
			notMain = true
		}
	case `windows`:
		nowPath := filepath.Base(os.Args[0])
		if nowPath[:7] == `___Test` || nowPath[len(nowPath)-9:] == `.test.exe` {
			notMain = true
		}
	}
	if notMain {
		_, onPath, _, _ := runtime.Caller(0)
		confPath = filepath.Join(onPath, `..`, `..`, `config.yaml`)
	}

	// 读取文件所有内容装到 []byte 中
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Println(`配置文件读取错误`)
		panic(err)
	}

	// 调用 Unmarshall 去解码文件内容
	err = yaml.Unmarshal(bytes, &C)
	if err != nil {
		log.Println(`配置文件解析错误`)
		panic(err)
	}

	// 如果有远程配置
	if C.Url != `` {
		var client = &http.Client{
			Timeout: 5 * time.Second, // 超时时间:5秒
			Transport: &http.Transport{ // 解决https证书
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		resp, err := client.Get(C.Url)
		if err != nil {
			log.Println(`获取远程配置错误 client.Get `, err)
			panic(err)
		}
		defer resp.Body.Close()

		// 读取返回值
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(`获取远程配置错误 ioutil.ReadAll `, err)
			panic(err)
		}

		var res struct {
			Code int    // 代码
			Data string // 数据集
			Msg  string // 消息
		}
		_ = json.Unmarshal(body, &res)

		if res.Code != 200 {
			log.Println(`获取远程配置错误 res.Code != 200 `, err)
			panic(res.Msg)
		}

		decryptCode := aes.Decrypt(res.Data, C.CKey)

		var t []Table
		_ = json.Unmarshal(decryptCode, &t)
		C.Tables = t
	}
}

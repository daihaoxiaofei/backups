package core

import (
	"backups/config"
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

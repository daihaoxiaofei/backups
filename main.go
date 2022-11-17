package main

import (
	"backups/core"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	core.Cron()
}

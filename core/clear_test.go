package core

import (
	"log"
	"testing"
)

func TestSaveDate(t *testing.T) {
	// log.Println(`2022/11/26 13:54:07 2022-11-26 13_00_00`[:13])

	for _, v := range saveDate() {
		log.Println(v)
	}

}

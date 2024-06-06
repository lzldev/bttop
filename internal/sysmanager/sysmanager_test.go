package sysmanager

import (
	"fmt"
	"testing"
	"time"
)

func TestSysManager(t *testing.T) {
	fmt.Println("yep")
	logger := make(chan string)

	go func() {
		for msg := range logger {
			fmt.Println(msg)
		}
	}()

	tx, rx := StartSysManager(logger)

	go func() {
		for msg := range rx {
			fmt.Printf("RX: %+v \n", msg)
		}
	}()

	tx <- SpeedUp
	time.Sleep(time.Second * 2)
	tx <- SpeedDown
	time.Sleep(time.Second * 2)

	time.Sleep(time.Second * 10)
}

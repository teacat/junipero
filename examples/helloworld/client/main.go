package main

import (
	"fmt"
	"time"

	"github.com/teacat/junipero"
)

func main() {
	c, _, err := junipero.NewClient(&junipero.ClientConfig{
		Address: "ws://127.0.0.1:8899/",
	})
	if err != nil {
		panic(err)
	}
	err = c.Write(time.Now().String())
	fmt.Printf("已發送訊息：%s\n", time.Now().String())
	if err != nil {
		panic(err)
	}
	msg, err := c.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("已接收訊息：%s\n", msg)
	err = c.Disconnect()
	if err != nil {
		panic(err)
	}
}

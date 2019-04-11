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
	go func() {
		for {
			err := c.Write(time.Now().String())
			fmt.Println("已發送訊息")
			if err != nil {
				panic(err)
			}
			<-time.After(time.Second * 1)
		}
	}()
	for {
		msg, err := c.Read()
		if err != nil {
			panic(err)
		}
		fmt.Printf("已接收訊息：%s\n", msg)
	}
}

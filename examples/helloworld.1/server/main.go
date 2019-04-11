package main

import (
	"log"
	"net/http"

	"github.com/teacat/junipero"
)

type server struct {
}

func (s *server) Close(sess *junipero.Session, status junipero.CloseStatus, msg string) error {
	log.Printf("連線被關閉：%d", status)
	return nil
}

func (s *server) Disconnect(sess *junipero.Session) {
	log.Println("連線正常關閉")
}

func (s *server) Error(sess *junipero.Session, err error) {
	log.Printf("連線錯誤：%s", err.Error())
}

func (s *server) Message(sess *junipero.Session, msg string) {
	log.Printf("接收到訊息：%s", msg)
	sess.Write(msg + " from server!")
}

func (s *server) MessageBinary(sess *junipero.Session, msg []byte) {
	log.Printf("接收到二進制訊息：%s", string(msg))
}

func (s *server) SentMessage(sess *junipero.Session, msg string) {
	log.Printf("已發送訊息：%s", msg)
}

func (s *server) SentMessageBinary(sess *junipero.Session, msg []byte) {
	log.Printf("已發送二進制訊息：%s", string(msg))
}

func (s *server) Pong(sess *junipero.Session) {
	log.Printf("接收到反響")
}

func (s *server) Ping(sess *junipero.Session) {
	log.Printf("接收到 Ping 請求")
}

func (s *server) Connect(sess *junipero.Session) {
	log.Printf("已連線")
}

func (s *server) Request(w http.ResponseWriter, r *http.Request, sess *junipero.Session) {
	log.Printf("有新的連線要升級至 WebSocket 通訊：%s", r.RemoteAddr)
}

func main() {
	j := junipero.NewServer(junipero.DefaultConfig(), &server{})
	http.HandleFunc("/", j.HandlerFunc())
	log.Fatal(http.ListenAndServe(":8899", nil))
}

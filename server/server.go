package server

import (
	"fmt"
	"im-system/user"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	Ip      string
	Port    int
	Users   map[string]*user.User
	Lock    sync.Mutex
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{Ip: ip, Port: port, Users: make(map[string]*user.User), Message: make(chan string)}
	return server
}
func (this *Server) handle(conn net.Conn) {
	// fmt.Println("连接成功")
	user := user.NewUser(conn)
	this.Lock.Lock()
	this.Users[user.Name] = user
	this.Lock.Unlock()
	this.Broadcast(user, "上线")
	select {}
}
func (this *Server) Broadcast(user *user.User, msg string) {
	sendMsg := fmt.Sprintf("%s : %s", user.Name, msg)
	this.Message <- sendMsg
}
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.Lock.Lock()
		for _, user := range this.Users {
			user.C <- msg
		}
		this.Lock.Unlock()
	}
}
func (this *Server) Start() {
	listener, err := net.Listen("tcp", this.Ip+":"+strconv.Itoa(this.Port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("错误：", err)
		}
	}()
	go this.ListenMessage()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go this.handle(conn)
	}

}

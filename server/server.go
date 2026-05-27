package server

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	ip   string
	port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{ip: ip, port: port}
	return server
}
func (this *Server) handle(conn net.Conn) {
	fmt.Println("连接成功")
}
func (this *Server) Start() {
	listener, err := net.Listen("tcp", this.ip+":"+strconv.Itoa(this.port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		go this.handle(conn)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("错误：", err)
		}
	}()
}

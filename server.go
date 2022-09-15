package main

import (
	"fmt"
	"net"
)

//定义Server结构体
type Server struct {
	Ip   string
	Port int
}

//创建一个Server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

//处理业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接已成功")
}

//开启Server
func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	fmt.Println(fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("net.Listen err :", err)
		return
	}
	defer listener.Close()
	//死循环处理业务
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err :", err)
			continue
		}
		//do something
		go server.Handler(conn)
	}

}

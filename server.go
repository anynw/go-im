package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

//定义Server结构体
type Server struct {
	Ip   string
	Port int
	//当前在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	//消息广播的channel
	Message chan string
}

//创建一个Server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//处理业务
func (this *Server) Handler(conn net.Conn) {
	// fmt.Println("链接已成功")
	user := NewUser(conn)
	//用户上线，将用户加入到OnlineMap中
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	//广播当前用户上线消息
	this.BroadCast(user, "上线了")

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				this.BroadCast(user, "已下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err : ", err)
			}

			msg := string(buf[:n-1])
			//广播
			this.BroadCast(user, msg)
		}
	}()

	//当前handler阻塞
	select {}
}

//广播消息方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

//监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线user
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		//将消息发送给全部的在线用户 【循环遍历 加锁】
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()

	}
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
	//启动监听Message的goroutine
	go server.ListenMessage()
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

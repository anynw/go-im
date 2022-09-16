package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
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
	user := NewUser(conn, this)

	user.Online()

	//监听用户是否活跃
	isLive := make(chan bool)

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				// this.BroadCast(user, "已下线")
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err : ", err)
			}

			msg := string(buf[:n-1])
			//广播
			// this.BroadCast(user, msg)
			user.DoMessage(msg)

			//用户发了消息，就代表活跃
			isLive <- true

		}
	}()

	//当前handler阻塞
	// select {}
	//添加超时强踢功能，使用定时器重置上线
	for {
		select {
		case <-isLive:
			//当前用户活跃，重置定时器
			//不做任何处理，为了激活select，更新下面的定时器
		case <-time.After(time.Second * 10):
			//说明超时了 将当前User强退
			// delete(this.OnlineMap, user)
			user.SendMsg("您已超时离线")
			//销毁用户资源
			close(user.C)
			//关闭链接
			conn.Close()
			//退出当前handler
			//return
			runtime.Goexit()
		}
	}
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

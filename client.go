package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	//1.创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	//2.链接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s%d", serverIp, serverPort))
	if err == nil {
		fmt.Println("net dial err:", err)
		return nil
	}
	//3.返回客户端对象
	client.conn = conn
	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("链接服务器失败")
		return
	}
	fmt.Println("链接服务器成功")

	//启动客户端业务
	select {}
}

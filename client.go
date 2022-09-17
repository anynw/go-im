package main

import (
	"flag"
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

//定义全局变量
var serverIp string
var serverPort int

//先于main函数调用
// .client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址，默认值127.0.0.1")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口，默认值8888")
}

func main() {
	//命令行解析
	flag.Parse()
	// client := NewClient("127.0.0.1", 8888)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("链接服务器失败")
		return
	}
	fmt.Println("链接服务器成功")

	//启动客户端业务
	select {}
}

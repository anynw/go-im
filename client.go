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
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//1.创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//2.链接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	fmt.Println("client conn sucess：", conn)
	if err != nil {
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

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}
}

func (client *Client) Run() {
	//循环调用客户端的flag 判断是否是退出模式
	for client.flag != 0 {
		//输入数字非法
		for client.menu() != true {

		}
		//根据输入的不同 处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式...")
			break
		case 2:
			//私聊模式
			fmt.Println("私聊模式...")
			break
		case 3:
			//更新用户名
			fmt.Println("更新用户名...")
			break

		}
	}
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
	// select {}
	client.Run()
}

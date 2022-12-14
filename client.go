package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
			// fmt.Println("公聊模式...")
			client.PublicChat()
			break
		case 2:
			//私聊模式
			// fmt.Println("私聊模式...")
			client.PrivateChat()
			break
		case 3:
			//更新用户名
			// fmt.Println("更新用户名...")
			client.UpdateName()
			break

		}
	}
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("请输入聊天内容,exit退出：")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("client conn err,", err)
				break
			}

		}

		chatMsg = ""
		fmt.Println("请输入聊天内容,exit退出：")
		fmt.Scanln(&chatMsg)

	}
}

func (client *Client) PrivateChat() {

	var remoteName string
	var chatMsg string

	client.selectOnlineUsers()
	fmt.Println("请输入聊天对象的用户名,exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("client conn err,", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("请输入聊天内容,exit退出：")
			fmt.Scanln(&chatMsg)
		}
		client.selectOnlineUsers()
		fmt.Println("请输入聊天对象的用户名,exit退出")
		fmt.Scanln(&remoteName)
	}

}

//查询在线用户
func (client *Client) selectOnlineUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client conn write err : ", err)
		return
	}

}

//更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名：")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client conn err:", err)
		return false
	}
	return true

}

//处理服务端goroutine 直接标准输出
func (client *Client) DealResponse() {
	//一旦conn有数据，直接拷贝到标准输出上，永久阻塞
	io.Copy(os.Stdout, client.conn)
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

	//单独开启goroutine处理server的回执消息
	go client.DealResponse()

	fmt.Println("链接服务器成功")

	//启动客户端业务
	// select {}
	client.Run()
}

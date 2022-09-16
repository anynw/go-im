package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
	//当前用户属于哪个server
	server *Server
}

//创建user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}
	go user.ListenMessage()
	return user
}

//用户上线业务
func (this *User) Online() {
	//用户上线，将用户加入到OnlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//广播当前用户上线消息
	this.server.BroadCast(this, "上线了")
}

//用户下线业务
func (this *User) Offline() {
	this.server.mapLock.Lock()
	//删除
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "下线了")
}

//用户处理消息业务
func (this *User) DoMessage(msg string) {
	//查询当前在线用户有哪些
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			// 查看非自己的登陆用户都有谁
			if this.Addr != user.Addr {
				onlineMsg := "[" + user.Addr + "]" + user.Name + ":在线......\n"
				this.SendMsg(onlineMsg)
			}
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式：rename|章三
		// newName := strings.Split(msg,"|")[1]
		newName := msg[7:]
		//判断当前登陆用户是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名被登陆\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("用户名已更新为" + this.Name + "\n")
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}

//给当前用户对应的客户端发送消息 【who 命令谁发起的，发给谁】
func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

//监听当前user channel的方法，一旦有消息，就直接发送给客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}

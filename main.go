/*
@Time : 2020/8/20 下午8:19
@Author : xiukang
@File : main.go
@Software: GoLand
*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go-socket/model"
	"go-socket/routers"
	"go-socket/runtime"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var mutex sync.Mutex

//// 用户退出
func logout(msg model.ReceiveMessage) {
	msg.Method = "logout"
	model.MessageBroadcast <- msg
}

// 消息处理
func onMessage(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("upgrade:", err)
		return
	}
	// 1.接收到用户连接,执行登录
	var loginMsg model.ReceiveMessage

	c.ReadJSON(&loginMsg)
	loginMsg.Client = c
	// 登录消息写入队列
	model.MessageBroadcast <- loginMsg

	// 关闭连接需要修改
	defer logout(loginMsg)
	for {
		// 1.处理当前用户获取系统消息
		var userMsg model.ReceiveMessage
		err := c.ReadJSON(&userMsg)
		userMsg.Client = c
		if err != nil {
			fmt.Print("read:", err)
			break
		}
		mes, _ := json.Marshal(userMsg)
		runtime.Info.Println("收到消息：->" + string(mes))
		model.MessageBroadcast <- userMsg
	}
}

/**
处理用户消息
*/
func handleMessages() {
	for {
		// 获取到管道里的所有数据
		msg := <-model.MessageBroadcast
		code, res := routers.GetRouter(msg)
		if code != 0 {
			runtime.Error.Println(res)
			msg.Message = "操作失败请重试"
			msg.Type = 2
			msg.Client.WriteJSON(msg)

			model.ErrorMessage = append(model.ErrorMessage, msg)
		}
	}
}

/**
处理错误的消息
*/
func handleErrorMessages() {
	for {
		for i, _ := range model.ErrorMessage { //range returns both the index and value
			// 错误的删除掉
			mutex.Lock()
			{
				runtime.Info.Println("处理错误消息")

				model.ErrorMessage = append(model.ErrorMessage[:i], model.ErrorMessage[i+1:]...)
			}
			// 释放锁，允许其他
			mutex.Unlock()
		}
		time.Sleep(time.Second * 1)
	}
}

// 代码初始化
func init() {
	runtime.Info.Println("系统初始化")
	// 1.修复用户会话与系统会话
	//sessionServer := &session.InitSession{}
	//go sessionServer.Index()
}

// 更新消息入mysql
func updateMessgeToDb() {
	//messageServer := &message.UpdateMessageQueueServer{}
	//for {
	//	messageServer.Index()
	//}

}

// 更新系统会话入mysql
func updateSessionToDb() {
	//messageServer := &session.UpdateSessionQueueServer{}
	//for {
	//	messageServer.Index()
	//}
}
func main() {
	go http.HandleFunc("/", onMessage)
	// 处理用户消息
	go handleMessages()
	//处理错误消息
	//go handleErrorMessages()
	//更新消息队列入mysql
	//go updateMessgeToDb()
	// 更新系统会话入mysql
	//go updateSessionToDb()

	runtime.Info.Println("websocket start at 127.0.0.1:3001")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

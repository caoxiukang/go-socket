package model

import "github.com/gorilla/websocket"

// 接收定义消息结构
type ReceiveMessage struct {
	// 请求的方法
	Method string `json:"method"`
	// 消息类型(1,text,0:系统消息,2:错误消息)
	Type int `json:"type"`
	// 消息体
	Message string `json:"message"`
	// 消息来源用户Id
	FromId uint64 `json:"fromId"`
	// 要发送用户Id
	ToId uint64 `json:"toId"`
	//// 当前连接
	Client *websocket.Conn `json:"client"`
}

//
// 所有用户连接辞池子
var Clients = make(map[uint64]*websocket.Conn) // connected clients

// 所有消息管道
var MessageBroadcast = make(chan ReceiveMessage) // broadcast channel

// 所有处理失败消息管道
var ErrorMessage = [] ReceiveMessage{}

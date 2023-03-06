package models

import (
	"encoding/json"
	"fatfoxChats/utils"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

// 消息关系
type Message struct {
	gorm.Model
	FormId   uint64 //发送者
	TargetId uint64 //接受者
	Type     int    //发送类型
	Media    string //文字图片音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //统计数字
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 存储映射关系
var clientMap = make(map[uint64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(w http.ResponseWriter, request *http.Request) {
	//  1.校验token 等合法性
	token := request.Header.Get("AUTHORIZATION")
	cal, err := utils.ParseToken(token)
	if err != nil {
		return
	}
	u := cal.Id // 取出url后面携带的参数userId
	fmt.Println(u)
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, request, nil)

	if err != nil {
		return
	}
	sId := strconv.FormatUint(u, 10)
	IsLoginState(sId, true)
	//	2.获取连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//	3.userid 和node绑定起来并且加锁
	rwLocker.Lock()
	clientMap[u] = node
	rwLocker.Unlock()
	//  4.完成发送的逻辑
	go sendProc(node, sId) //第二步执行
	//	5.完成接受逻辑
	go recvProc(node, sId) //第一步执行
	//sendMsg(userId, []byte("欢迎进入聊天室"))
}

// 写入到链接读取之后的数据
func sendProc(node *Node, userId string) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				IsLoginState(userId, false)
				fmt.Println(err)
				return
			}
		}
	}
}

// 循环，读取链接中的数据
func recvProc(node *Node, userId string) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			IsLoginState(userId, false)
			return
		}

		type M struct {
			FormId   uint64
			TargetId uint64
			Type     int
			Content  string
			SendTime int64
		}
		//	3.校验用户关系
		var m M
		err = json.Unmarshal([]byte(data), &m)

		if err != nil {
			return
		}
		FId := strconv.FormatUint(m.FormId, 10)
		TId := strconv.FormatUint(m.TargetId, 10)
		b := IsFriends(FId, TId)
		if b {
			broadMsg(data)
		} else {
			id, e := strconv.ParseUint(FId, 10, 64)
			if e != nil {
				return
			}
			sendMsg(id, []byte("{\"code\":1001,\"msg\":\"对不起，该用户不是你好友！！！\"}"))
		}
	}
}

var udpsendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp数据发送协程
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		// 如果udpsendChan
		case data := <-udpsendChan:
			_, errs := conn.Write(data)
			if errs != nil {
				fmt.Println(errs)
				return
			}

		}
	}
}

// 完成udp数据接受协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		var buf [512]byte
		n, errs := con.Read(buf[0:])
		if errs != nil {
			fmt.Println(errs)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	errs := json.Unmarshal(data, &msg)
	if errs != nil {
		fmt.Println(errs)
		return
	}
	//客户端发送消息的类型
	switch msg.Type {
	case 1:
		sendMsg(msg.TargetId, data)
	case 2:
		//私信
	case 3:
		//群发
	}
}
func sendMsg(TargetId uint64, msg []byte) {
	rwLocker.Lock()
	node, ok := clientMap[TargetId]
	rwLocker.Unlock()
	if ok {
		node.DataQueue <- msg
	}
}

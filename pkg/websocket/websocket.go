package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/tiptok/gocomm/pkg/mybeego"
	"reflect"
	"sync"
)

type ConnState int

const (
	Disconnected ConnState = iota
	Connected
)

func init() {
	keyType := reflect.TypeOf(&websocket.Conn{})
	valueType := reflect.TypeOf(&WebsocketConnection{})
	Connections = NewJMap(keyType, valueType)
	Clients = NewJMap(reflect.TypeOf("1:1"), valueType)
}

type ReceiveHandler (func([]byte) *mybeego.Message)

type WebsocketConnection struct {
	Uid         int64
	AppId       int
	Conn        *websocket.Conn
	Echan       chan interface{}
	Wchan       chan string
	State       ConnState
	OnReceive ReceiveHandler
	OnceClose  sync.Once
}

func NewWebsocketConnection(conn *websocket.Conn,head *mybeego.RequestHead,recv ReceiveHandler)*WebsocketConnection{
	return &WebsocketConnection{
		Uid:         head.Uid,
		AppId:       head.AppId,
		Conn:        conn,
		Echan:       make(chan interface{}),
		Wchan:       make(chan string, 10),
		State:       Connected,
		OnReceive: recv,
	}
}

//声明了两个cliets 管理 一个通过uid  一个通过conn管理
// key(*websocket.Conn) value(*WebsocketConnection)
var Connections *JMap
// key=uid(int64) value(*WebsocketConnection)
var Clients *JMap

type JMap struct {
	sync.RWMutex
	m         map[interface{}]interface{}
	keyType   reflect.Type
	valueType reflect.Type
}

func NewJMap(keyType, valueType reflect.Type) *JMap {
	return &JMap{
		keyType:   keyType,
		valueType: valueType,
		m:         make(map[interface{}]interface{}),
	}
}

func (this *JMap) PrintConnectStatus() interface{} {
	beego.Debug("PrintConnectStatus...")
	beego.Info("============查看websocket连接状态begin============")
	for i, v := range this.m {
		beego.Info("key:", i, " conn:", v)
	}
	beego.Info("============查看websocket连接状态end============")
	return this.m
}

func (this *JMap) GetOnlineClient() map[interface{}]interface{} {
	return this.m
}

func (this *JMap) acceptable(k, v interface{}) bool {
	if k == nil || reflect.TypeOf(k) != this.keyType {
		return false
	}

	if k == nil || reflect.TypeOf(v) != this.valueType {
		return false
	}

	return true
}

func (this *JMap) Get(k interface{}) (interface{}, bool) {
	this.RLock()
	conn, ok := this.m[k]
	this.RUnlock()
	return conn, ok
}

func (this *JMap) Put(k interface{}, v interface{}) bool {
	if !this.acceptable(k, v) {
		return false
	}
	if connI, ok := Clients.Get(k); ok {
		beego.Debug("key:", k, "已经连接,先剔除下线")
		if conn, ok := connI.(*WebsocketConnection); ok {
			//conn.Conn.WriteMessage(websocket.TextMessage, []byte("您的帐号在其它地方登录,您被剔除下线"))
			conn.Close()
		}

	}
	this.Lock()
	this.m[k] = v
	this.Unlock()
	return true
}

func (this *JMap) Remove(k interface{}) {
	this.Lock()
	delete(this.m, k)
	this.Unlock()
}

func (this *JMap) Clear() {
	this.Lock()
	this.m = make(map[interface{}]interface{})
	this.Unlock()
}

func (this *JMap) Size() int {
	this.RLock()
	defer this.RUnlock()
	return len(this.m)
}

func (this *JMap) IsEmpty() bool {
	return this.Size() == 0
}

func (this *JMap) Contains(k interface{}) bool {
	this.RLock()
	_, ok := this.m[k]
	this.RUnlock()
	return ok
}

func (c *WebsocketConnection) Serve() {
	c.State = Connected
	Connections.Put(c.Conn, c)
	key := fmt.Sprintf("%d:%d", c.Uid, c.AppId)
	Clients.Put(key, c)

	go doWrite(c)
	doRead(c)
}

func (c *WebsocketConnection) Send(msg string) {
	//panic("panic in websocket.send...")
	c.Wchan <- msg
}

func (c *WebsocketConnection) Close() {
    c.OnceClose.Do(func(){
		beego.Info("ws:close----uid:", c.Uid, "appid:", c.AppId, "state:", c.State)
		if c.State == Disconnected {
			return
		}

		Connections.Remove(c.Conn)
		key := fmt.Sprintf("%d:%d", c.Uid, c.AppId)
		Clients.Remove(key)

		c.State = Disconnected
		close(c.Echan)
		close(c.Wchan)
		c.Conn.Close()
	})
}

func doRead(c *WebsocketConnection) {
	defer func() {
		//beego.Debug("doRead exit...uid:", c.Uid, "appid:", c.AppId)
		c.Close()
	}()

	for {
		select {
		case <-c.Echan:
			return
		default:
		}
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			beego.Info(err)
			return
		}
		beego.Info(fmt.Sprintf("===>ws:recv msg from uid(%d) : %s", c.Uid, string(msg)))
		retMsg := c.OnReceive(msg)
		retMsgByte, err := json.Marshal(retMsg)
		beego.Info(fmt.Sprintf("<===ws:send to client uid(%d) : %s", c.Uid, string(retMsgByte)))
		c.Send(string(retMsgByte))
	}
}

func doWrite(c *WebsocketConnection) {
	defer func() {
		if err := recover(); err != nil {
			beego.Error("Recover in doWrite...uid:", c.Uid, "apid:", c.AppId, "err:", err)
		}
	}()
	defer func() {
		//beego.Debug("doWrite exit...uid:", c.Uid, "appid:", c.AppId)
		c.Close()
	}()

	for {
		select {
		case <-c.Echan:
			return
		default:
		}
		msg, ok := <-c.Wchan
		if !ok {
			break
		}
		err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			break
		}
	}
}

func SendDataByWs(uid int64, appId int, sendMsg interface{}) bool {
	if sendMsg == nil || uid < 1 || appId < 1 {
		return false
	}
	msg := &mybeego.Message{
		Errno:  0,
		Errmsg: mybeego.NewMessage(0).Errmsg,
		Data:   sendMsg,
	}
	msgByte, err := json.Marshal(msg)
	if err != nil {
		beego.Error(err)
		return false
	}
	key := fmt.Sprintf("%d:%d", uid, appId)
	if connI, ok := Clients.Get(key); ok {
		beego.Debug(ok)
		if conn, ok := connI.(*WebsocketConnection); ok {
			conn.Send(string(msgByte))
			return true
		}
	}
	return false
}

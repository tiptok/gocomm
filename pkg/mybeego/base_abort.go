package mybeego

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/tiptok/gocomm/pkg/log"
//	"github.com/tiptok/gocomm/xtime"
//	"strconv"
//	"time"
//
//	"github.com/astaxie/beego"
//	"github.com/tiptok/gocomm/pkg/redis"
//)
//
//// BaseController
//type BaseController struct {
//	beego.Controller
//	Query       map[string]string
//	JSONBody    map[string]interface{}
//	ByteBody    []byte
//	RequestHead *RequestHead
//}
//
//func assertCompleteImplement() {
//	var _ beego.ControllerInterface = (*BaseController)(nil)
//}
//
//func (this *BaseController) Options() {
//	this.AllowCross() //允许跨域
//	this.Data["json"] = map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}
//	this.ServeJSON()
//}
//
//func (this *BaseController) AllowCross() {
//	//c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")       //允许访问源
//	//c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
//	//this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
//	//c.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
//	//c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
//	//c.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
//
//	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://cdn.didong123.cn")
//	//this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
//	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "uid, token,jwt, deviceid, appid,Content-Type,Authorization,from")
//	this.Ctx.WriteString("")
//
//}
//
//func (this *BaseController) Prepare() {
//
//	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
//	if this.Ctx.Input.Method() == "OPTIONS" {
//		this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//		this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "uid, token,jwt, deviceid, appid,Content-Type,Authorization,from")
//		this.Ctx.WriteString("")
//		return
//	}
//
//	this.Query = map[string]string{}
//	input := this.Input()
//	for k := range input {
//		this.Query[k] = input.Get(k)
//	}
//	if this.Ctx.Input.RequestBody != nil {
//		// contentType := this.Ctx.Input.Header("Content-type")
//		// if strings.HasPrefix(contentType, "application/json") {
//		this.ByteBody = this.Ctx.Input.RequestBody[:]
//		if len(this.ByteBody) < 1 {
//			this.ByteBody = []byte("{}")
//		}
//		this.RequestHead = &RequestHead{}
//		this.RequestHead.Token = this.Ctx.Input.Header("token")
//		this.RequestHead.Version = this.Ctx.Input.Header("version")
//		this.RequestHead.Os = this.Ctx.Input.Header("os")
//		this.RequestHead.From = this.Ctx.Input.Header("from")
//		this.RequestHead.Screen = this.Ctx.Input.Header("screen")
//		this.RequestHead.Model = this.Ctx.Input.Header("model")
//		this.RequestHead.Channel = this.Ctx.Input.Header("channel")
//		this.RequestHead.Net = this.Ctx.Input.Header("net")
//		this.RequestHead.DeviceId = this.Ctx.Input.Header("deviceid")
//		this.RequestHead.Uid, _ = strconv.ParseInt(this.Ctx.Input.Header("uid"), 10, 64)
//		this.RequestHead.AppId, _ = strconv.Atoi(this.Ctx.Input.Header("appid"))
//		this.RequestHead.LoginIp = this.Ctx.Input.IP()
//		this.RequestHead.Jwt = this.Ctx.Input.Header("jwt")
//		this.RequestHead.SetRequestId(fmt.Sprintf("%v.%v.%s", this.RequestHead.Uid, time.Now().Format(xtime.YYYYMMDDHHMMSS), this.Ctx.Request.URL))
//		log.Info(fmt.Sprintf("====>Recv data from uid(%d) client:\nHeadData: %s\nRequestId:%s BodyData: %s", this.RequestHead.Uid, this.Ctx.Request.Header, this.RequestHead.GetRequestId(), string(this.ByteBody)))
//	}
//	key := SWITCH_INFO_KEY
//	str := ""
//	switchInfo := &TotalSwitchStr{}
//	if str, _ = redis.Get(key); str == "" {
//		switchInfo.TotalSwitch = TOTAL_SWITCH_ON
//		switchInfo.MessageBody = "正常运行"
//		redis.Set(key, switchInfo, redis.INFINITE)
//	} else {
//		json.Unmarshal([]byte(str), switchInfo)
//	}
//	if switchInfo.TotalSwitch == TOTAL_SWITCH_OFF {
//		var msg *Message
//		msg = NewMessage(3)
//		msg.Errmsg = switchInfo.MessageBody
//		log.Info(msg.Errmsg)
//		this.Data["json"] = msg
//		this.ServeJSON()
//		return
//	}
//}
//
//func (this *BaseController) Resp(msg *Message) {
//
//	this.Data["json"] = msg
//	this.ServeJSON()
//}
//
//func (this *BaseController) Finish() {
//
//	if this.Ctx.Input.Method() == "OPTIONS" {
//		return
//	}
//
//	strByte, _ := json.Marshal(this.Data["json"])
//	length := len(strByte)
//	if length > 5000 {
//		log.Info(fmt.Sprintf("<====Send to uid(%d) client: %d byte\nRequestId:%s RspBodyData: %s......", this.RequestHead.Uid, length, this.RequestHead.GetRequestId(), string(strByte[:5000])))
//	} else {
//		log.Info(fmt.Sprintf("<====Send to uid(%d) client: %d byte\nRequestId:%s RspBodyData: %s", this.RequestHead.Uid, length, this.RequestHead.GetRequestId(), string(strByte)))
//	}
//}
//
//// BaseControllerCallBack
//type BaseControllerCallBack struct {
//	beego.Controller
//	Query    map[string]string
//	JSONBody map[string]interface{}
//	ByteBody []byte
//}
//
//func (this *BaseControllerCallBack) Prepare() {
//	this.Query = map[string]string{}
//	input := this.Input()
//	for k := range input {
//		this.Query[k] = input.Get(k)
//	}
//
//	if this.Ctx.Input.RequestBody != nil {
//		log.Info("RecvHead:", string(this.Ctx.Input.Header("Authorization")))
//		this.ByteBody = this.Ctx.Input.RequestBody
//	}
//}
//
//func (this *BaseControllerCallBack) Resp(msg *Message) {
//	this.Data["json"] = msg
//	this.ServeJSON()
//}
//
//func (this *BaseControllerCallBack) Finish() {
//	strByte, _ := json.Marshal(this.Data["json"])
//	log.Debug("<====Send to client:\n", string(strByte))
//}

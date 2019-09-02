package redis

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"sync"
)

type Mutex struct {
	conn redis.Conn
	timeOut int64
	resource string
	lock bool
	closeOnce sync.Once
}
//NewMutex  create new mutex
func NewMutex(source string)*Mutex{
	return &Mutex{
		resource:source,
		lock:false,
		timeOut:SECOND*5,//未执行完,已经超时 超时时间设大
	}
}
func (l *Mutex)Key()string{
	return fmt.Sprintf("reidslock:%s",l.resource)
}
func(l *Mutex)Conn()redis.Conn{
	return l.conn
}
//设置超时
func (l *Mutex)TimeOut(t int64)*Mutex{
	l.timeOut = t
	return l
}
//加锁
//true:加锁成功  false:加锁失败
func (l *Mutex)Lock()( bool){
	defer func(){
		if !l.lock{
			beego.Warn("on locked:",l.Key())
			l.Close()
		}
	}()
	if l.lock{
		return l.lock
	}
	l.conn =RedisClient.Get()
	resourceKey :=l.Key()
	if result, err := l.conn.Do("SET", resourceKey,l.resource,"NX","EX", l.timeOut); err != nil || result==nil{
		return l.lock
	}else{
		ok :=result.(string)
		if ok !="OK"{
			return l.lock
		}
	}
	l.lock = true
	return l.lock
}
//解锁
func(l *Mutex)UnLock()error{
	defer l.Close()
	if !l.lock{
		return nil
	}
	if _,err :=l.conn.Do("DEL",l.Key());err!=nil{
		return err
	}
	l.lock = false
	return nil
}
//关闭
func(l *Mutex)Close(){
	l.closeOnce.Do(func(){
		if l.conn!=nil{
			l.conn.Close()
		}
	})
}

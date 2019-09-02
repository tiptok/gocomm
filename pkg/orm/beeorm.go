package orm

import (
	"github.com/astaxie/beego/orm"
	"github.com/tiptok/gocomm/config"
	"github.com/tiptok/gocomm/pkg/log"
)

func NewBeeormEngine(conf config.Mysql){
	err:=orm.RegisterDataBase("default","mysql",conf.DataSource)
	if err!=nil{
		log.Error(err)
	}else{
		log.Debug("open db address:",conf.DataSource)
	}
	orm.SetMaxIdleConns("default", conf.MaxIdle)
	orm.SetMaxOpenConns("default", conf.MaxOpen)
}

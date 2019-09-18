package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type IConfig interface {
	Int(k string)(int,error)
	Int64(k string)(int64,error)

	String(k string)string
	Strings(k string)[]string

	Bool(k string)(bool,error)
	Float(k string)(float64,error)
}

func assertViperImplementIconfig(){
	var _ IConfig= (*ViperConfig)(nil)
}

type ViperConfig struct {
	*viper.Viper
}

func(v *ViperConfig)Int(k string)(value int,err error){
	defer func(){
		if p:=recover();p!=nil{
			err = fmt.Errorf("%v",p)
			return
		}
	}()
	value = v.GetInt(k)
	return
}
func(v *ViperConfig)Int64(k string)(value int64,err error){
	defer func(){
		if p:=recover();p!=nil{
			err = fmt.Errorf("%v",p)
			return
		}
	}()
	value = v.GetInt64(k)
	return 0,nil
}
func(v *ViperConfig)String(k string)string{
	return v.String(k)
}
func(v *ViperConfig)Strings(k string)[]string{
	return v.Strings(k)
}
func(v *ViperConfig)Bool(k string)(bool,error){
	return v.Bool(k)
}
func(v *ViperConfig)Float(k string)(float64,error){
	return v.Float(k)
}


var Default IConfig

func NewViperConfig(configType,configFile string){
	v := viper.New()
	v.SetConfigType(configType)
	v.SetConfigFile(configFile)
	Default = &ViperConfig{v}
}



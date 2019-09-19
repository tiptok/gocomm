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
	viper *viper.Viper
}

func(v *ViperConfig)Int(k string)(value int,err error){
	defer func(){
		if p:=recover();p!=nil{
			err = fmt.Errorf("%v",p)
			return
		}
	}()
	value = v.viper.GetInt(k)
	return
}
func(v *ViperConfig)Int64(k string)(value int64,err error){
	defer func(){
		if p:=recover();p!=nil{
			err = fmt.Errorf("%v",p)
			return
		}
	}()
	value = v.viper.GetInt64(k)
	return 0,nil
}
func(v *ViperConfig)String(k string)string{
	return v.viper.GetString(k)
}
func(v *ViperConfig)Strings(k string)[]string{
	return v.viper.GetStringSlice(k)
}
func(v *ViperConfig)Bool(k string)(b bool,err error){
	defer func(){
		if p:=recover();p!=nil{
			err = fmt.Errorf("%v",p)
			return
		}
	}()
	b = v.viper.GetBool(k)
	return
}
func(v *ViperConfig)Float(k string)(f float64,err error){
	defer func(){
		if p:=recover();p!=nil{
			err = fmt.Errorf("%v",p)
			return
		}
	}()
	f = v.viper.GetFloat64(k)
	return
}


var Default IConfig

func NewViperConfig(configType,configFile string)IConfig{
	v := viper.New()
	v.SetConfigType(configType)
	v.SetConfigFile(configFile)
	v.ReadInConfig()
	Default = &ViperConfig{v}
	return Default
}



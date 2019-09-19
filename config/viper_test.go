package config

import (
	"testing"
)

func Test_NewViperConfig(t *testing.T){
	NewViperConfig("yaml","F:\\examples_gincomm\\conf\\app-dev.yaml")
	dataSource :=Default.String("redis_url")
	if len(dataSource)==0{
		t.Fatal("error get")
	}
}

func Benchmark_NewViperConfig(b *testing.B){
	NewViperConfig("yaml","F:\\examples_gincomm\\conf\\app-dev.yaml")
	dataSource :=""
	for i:=0;i<b.N;i++{
		dataSource =Default.String("redis_url")
		if len(dataSource)==0{
			b.Fatal("error get")
		}
	}
}

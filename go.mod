module github.com/tiptok/gocomm

go 1.12

require (
	github.com/astaxie/beego v1.10.0

	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/garyburd/redigo v1.6.0
	github.com/gin-gonic/gin v1.4.0
	github.com/google/go-cmp v0.2.0
	github.com/gorilla/websocket v1.4.1
	github.com/mattn/go-sqlite3 v1.11.0 // indirect

	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9
	github.com/opentracing/opentracing-go v1.1.0
	github.com/spf13/viper v1.4.0
	github.com/uber/jaeger-client-go v2.16.0+incompatible
	github.com/uber/jaeger-lib v2.0.0+incompatible
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43

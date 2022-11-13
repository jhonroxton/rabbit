package main

import (
    "fmt"
	gl "rabbit/globals"
	"rabbit/mysql/db"
	"rabbit/mysql/inits"
	redis "rabbit/redis"
	"rabbit/routers"
	"os"
)

func init() {
	//初始化配置文件
	gl.V = inits.Viper()
	//初始化 redis
	builder := redis.NewRedisBuilder(false, gl.V.GetString("redis.addr")+":"+gl.V.GetString("redis.port"),
		gl.V.GetString("redis.password"), gl.V.GetInt("redis.db"))
	redisPool := redis.NewRedisPool().SetBuilder(builder)
	err := redisPool.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//初始化数据库
	db.Instance().GetDb()
}

func main() {
	//初始化gin路由
	routers.RunServive()

}
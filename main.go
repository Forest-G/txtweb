package main

import (
	"fmt"
	"net/http"
	"os"
	"txt/interfaces"
	"txt/pkg/config"
	"txt/redisdata"
	"txt/sqldata"

	"github.com/sirupsen/logrus"
)

func init() {
	f, err := os.OpenFile("logs/logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetReportCaller(true)
	logrus.SetOutput(f)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		//以下设置只是为了使输出更美观
		// DisableColors:   true,
		// TimestampFormat: "2006-01-02 15:04:05",

		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})
}

func main() {
	config.Init()          //读取配置文件
	redisdata.Init()       //连接redis
	sqldata.Opengowebsql() //打开数据库
	interfaces.Init()      //连接处理器函数
	fmt.Println("服务开启")
	if err := http.ListenAndServe(config.C.Port, nil); err != nil {
		logrus.Panic("监听失败")
	}
}

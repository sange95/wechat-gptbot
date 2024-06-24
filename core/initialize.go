package core

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"wechat-gptbot/config"
	"wechat-gptbot/core/handler"
	"wechat-gptbot/core/plugins"
	"wechat-gptbot/core/plugins/news"
	"wechat-gptbot/core/plugins/weather"
	"wechat-gptbot/core/svc"
	"wechat-gptbot/logger"
	"wechat-gptbot/streamlit_app"
)

var Bot *openwechat.Bot

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/8 17:12
* @Package:
 */
func Initialize() {
	// 初始化日志
	logger.InitLogrus(logger.Config{
		Level:      logrus.DebugLevel,
		ObjectName: "wechat-gptbot",
		WriteFile:  false,
	})

	Bot = openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 初始化插件
	plugins.Manger.Register(weather.NewPlugin(), news.NewPlugin())
	// 初始化配置文件
	config.InitConfig()
	// 初始化会话上下文管理器
	handler.Context = svc.NewServiceContext()

	// 启动streamlit
	go streamlit_app.RunStreamlit()

}

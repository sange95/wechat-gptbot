package main

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"wechat-gptbot/core"
	"wechat-gptbot/core/handler"
	"wechat-gptbot/server"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/8 10:46
* @Package:
 */
func main() {
	// 初始化核心配置
	core.Initialize()

	// 启动监听端口
	go server.NewApiServer(core.Bot).Run()
	// 定义消息处理函数
	// 获取消息处理器
	dispatcher := handler.NewMessageMatchDispatcher()
	core.Bot.MessageHandler = dispatcher.AsMessageHandler()
	core.Bot.UUIDCallback = consoleQrCode // 注册登陆二维码回调
	// 登录回调
	//bot.SyncCheckCallback = nil
	reloadStorage := openwechat.NewFileHotReloadStorage("token.json")
	if err := core.Bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); nil != err {
		panic(err)
	}
	// 获取当前登录的用户
	self, err := core.Bot.GetCurrentUser()
	if nil != err {
		panic(err)
	}
	go handler.KeepAlive(self)
	logrus.Infof("login success %s,%s,%s", self.User, self.City, self.Province)
	core.Bot.Block()
}

func consoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Medium)
	fmt.Println(q.ToSmallString(false))
}

package svc

import (
	"wechat-gptbot/config"
	"wechat-gptbot/core/chat_lm"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/12 11:09
* @Package:
 */

type ServiceContext struct {
	Session chat_lm.Session
}

func NewServiceContext() *ServiceContext {
	if config.C.BaseModel == "baidu" {
		return &ServiceContext{Session: chat_lm.NewBaiduSession()}
	}

	return &ServiceContext{
		Session: chat_lm.NewSession(),
	}
}

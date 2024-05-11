package main

import (
	"fmt"
	"wechat-gptbot/core/plugins/wechatMovement/zeepLife"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/5/10 14:59
* @Package:
 */
func main() {
	app := zeepLife.NewZeppLife("1003941268@knownsec.com", "4f4ezha!")
	err := app.SetSteps(7500)
	if err != nil {
		panic(err)
	}
	fmt.Println("success set step")
}

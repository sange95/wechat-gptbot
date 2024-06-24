package controller

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"wechat-gptbot/core"
)

var FriendsMap = make(map[string]*openwechat.Friend)
var GroupsMap = make(map[string]*openwechat.Group)

var once sync.Once

func Init() {
	self, _ := core.Bot.GetCurrentUser()
	once.Do(func() {
		groups, _ := self.Groups(false)
		friends, _ := self.Friends(false)

		for i, group := range groups {
			fmt.Println("UserName:", group.UserName)
			fmt.Println("NickName:", group.NickName)
			fmt.Println("DisplayName:", group.DisplayName)
			fmt.Println("DisplayName:", group.RemarkName)
			GroupsMap[group.UserName] = groups[i]
		}

		for i, friend := range friends {
			fmt.Println("UserName:", friend.UserName)
			fmt.Println("NickName:", friend.NickName)
			fmt.Println("DisplayName:", friend.DisplayName)
			fmt.Println("DisplayName:", friend.RemarkName)
			FriendsMap[friend.UserName] = friends[i]
		}
	})
}

type MsgEntity struct {
	Msg          string `json:"msg"`
	ReceiverType string `json:"receiverType"` // 群组，个人
	ReceiverName string `json:"receiverName"` // 接受者的名字，群名或者微信名称
}

func SendMsg(c *gin.Context) {

	param := MsgEntity{}
	Init()

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, NewFailureResponse(fmt.Sprintf("unmarsha param failure:%+v", err)))
		return
	}

	if param.Msg == "" {
		c.JSON(http.StatusOK, NewSuccessResponse("send success."))
		return
	}

	self, _ := core.Bot.GetCurrentUser()
	fmt.Printf("self:%+v", self)
	switch param.ReceiverType {
	case "friend":
		f := FindFriend(param.ReceiverName)
		if f == nil {
			c.JSON(http.StatusNotFound, NewFailureResponse(fmt.Sprintf("%s not found", param.ReceiverName)))
			return
		}
		if _, err := self.SendTextToFriend(f, param.Msg); err != nil {
			c.JSON(http.StatusNotFound, NewFailureResponse(fmt.Sprintf("send failed:%+v", err)))
			return
		}

	case "group":
		g := FindGroup(param.ReceiverName)
		if g == nil {
			c.JSON(http.StatusNotFound, NewFailureResponse(fmt.Sprintf("%s not found", param.ReceiverName)))
			return
		}
		if _, err := self.SendTextToGroup(g, param.Msg); err != nil {
			c.JSON(http.StatusNotFound, NewFailureResponse(fmt.Sprintf("send failed:%+v", err)))
			return
		}
	default:
		c.JSON(http.StatusBadRequest, NewFailureResponse("receiver type not support."))
		return
	}
	c.JSON(http.StatusOK, NewSuccessResponse("send success."))
	return
}

func FindFriend(username string) *openwechat.Friend {
	return FriendsMap[username]
}

func FindGroup(groupName string) *openwechat.Group {
	return GroupsMap[groupName]
}

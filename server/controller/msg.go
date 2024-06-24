package controller

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/gin-gonic/gin"
	"net/http"
	"wechat-gptbot/core"
)

var FriendsMap = make(map[string]*openwechat.Friend)
var GroupsMap = make(map[string]*openwechat.Group)

type MsgEntity struct {
	Msg          string `json:"msg"`
	ReceiverType string `json:"msgType"`      // 群组，个人
	ReceiverName string `json:"receiverName"` // 接受者的名字，群名或者微信名称
}

func SendMsg(c *gin.Context) {

	param := MsgEntity{}

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, NewFailureResponse(fmt.Sprintf("unmarsha param failure:%+v", err)))
		return
	}

	if param.Msg == "" {
		c.JSON(http.StatusOK, NewSuccessResponse("send success."))
		return
	}

	self, _ := core.Bot.GetCurrentUser()
	switch param.ReceiverType {
	case "friend":
		f := FindFriend(param.ReceiverName, self)
		if f == nil {
			c.JSON(http.StatusNotFound, NewFailureResponse(fmt.Sprintf("%s not found", param.ReceiverName)))
			return
		}
		if _, err := self.SendTextToFriend(f, param.Msg); err != nil {
			c.JSON(http.StatusNotFound, NewFailureResponse(fmt.Sprintf("send failed:%+v", err)))
			return
		}

	case "group":
		g := FindGroup(param.ReceiverName, self)
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

func FindFriend(username string, self *openwechat.Self) *openwechat.Friend {
	f, ok := FriendsMap[username]
	friends, _ := self.Friends(false)
	if !ok {
		for i, friend := range friends {
			if friend.UserName == username {
				f = friends[i]
			}
		}
	}
	return f
}

func FindGroup(groupName string, self *openwechat.Self) *openwechat.Group {
	g, ok := GroupsMap[groupName]
	groups, _ := self.Groups(false)
	if !ok {
		for i, group := range groups {
			if group.NickName == groupName {
				g = groups[i]
			}
		}
	}
	return g
}

package chat_lm

import (
	"context"
	"fmt"
	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/sashabaranov/go-openai"
	"sync"
	"time"
	"wechat-gptbot/config"
	"wechat-gptbot/core/plugins"
)

type baiduSpeedSession struct {
	Session
	sync.RWMutex                                 // 用户的创建需要加锁
	client         *qianfan.ChatCompletion       // 会话客户端
	ctx            map[string]*baiduUserMessage  // 管理用户上下文,根据用户 id 来个例
	prompt         qianfan.ChatCompletionMessage // 管理提示词
	pluginRegistry *plugins.PluginManger         // 插件注册器
	//image qianfan.Text2Image
}

func initBaiduPrompt() qianfan.ChatCompletionMessage {
	// 获取所有插件信息
	pluginsInfo := plugins.Manger.PluginPrompt()
	prompt := fmt.Sprintf(config.Prompt, time.Now().Format(time.DateTime), pluginsInfo)
	return qianfan.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}
}

func NewBaiduSession() Session {
	qianfan.GetConfig().AccessKey = ""
	qianfan.GetConfig().SecretKey = ""

	// 调用对话Chat，可以通过 WithModel 指定模型，例如指定ERNIE-3.5-8K，参数对应ERNIE-Bot
	client := qianfan.NewChatCompletion(
		qianfan.WithModel("ERNIE-Speed-128K"),
	)
	registry := plugins.NewPluginRegistry()
	return &baiduSpeedSession{
		RWMutex:        sync.RWMutex{},
		ctx:            make(map[string]*baiduUserMessage),
		client:         client,
		pluginRegistry: registry,
		prompt:         initBaiduPrompt(),
	}
}

// 获取用户
func (s *baiduSpeedSession) getUserContext(userName string) *baiduUserMessage {

	if msg, ok := s.ctx[userName]; ok {
		return msg
	}
	s.Lock()
	defer s.Unlock()
	// 双检加锁，防止加锁的过程中已经创建了用户
	if msg, ok := s.ctx[userName]; ok {
		return msg
	}
	msg := s.newUserMessage(userName)
	return msg
}

// Prompt 获取提示词  todo:将插件写入提示词
func (s *baiduSpeedSession) Prompt() qianfan.ChatCompletionMessage {
	return s.prompt
}

// 用户级消息
type baiduUserMessage struct {
	sync.Mutex                                 // 加锁 防止上下文顺序紊乱 一个用户只能拿到响应后才能再次提问
	user       string                          // 用户
	ctx        []qianfan.ChatCompletionMessage // 用户聊天的上下文 最多只保留6条记录，3组对话
}

// 新建一个用户级消息
func (s *baiduSpeedSession) newUserMessage(user string) *baiduUserMessage {
	msg := &baiduUserMessage{
		user:  user,
		ctx:   []qianfan.ChatCompletionMessage{},
		Mutex: sync.Mutex{},
	}
	s.ctx[user] = msg
	return msg
}

// 给用户追加上下文
func (um *baiduUserMessage) addContext(currentMessage, prompt qianfan.ChatCompletionMessage) {
	um.ctx = append(um.ctx, currentMessage)
	// 最多保存10条上下文
	if len(um.ctx) > MaxSession {
		um.ctx = um.ctx[len(um.ctx)-MaxSession:]
		// 将prompt 作为上下文第一条
		um.ctx[0] = prompt
		um.ctx[1] = qianfan.ChatCompletionAssistantMessage("好的，我会尽力办好群助理的角色，请问我有什么需要我帮助的吗？")
	}

}

// 构建上下文到消息体
func (um *baiduUserMessage) buildMessage(userName string, currentMsg qianfan.ChatCompletionMessage) []qianfan.ChatCompletionMessage {
	msgs := append(um.ctx, currentMsg)
	fmt.Println("=====" + userName + "=======")
	for i, ctx := range msgs {
		fmt.Printf("%d     %s\n", i, ctx.Content)
	}
	fmt.Println("=====" + userName + "=======")
	return msgs
}

func (s *baiduSpeedSession) Chat(ctx context.Context, content string) string {
	currentMsg := qianfan.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}
	//fmt.Println("==============")
	//fmt.Println(currentMsg.Role)
	// 默认不带上下文
	sender := ctx.Value("sender").(string)
	// 获取用户上下文
	um := s.getUserContext(sender)
	fmt.Println(len(um.ctx))
	var msgs []qianfan.ChatCompletionMessage
	if len(um.ctx) == 0 {
		um.addContext(s.Prompt(), s.Prompt())
		um.addContext(qianfan.ChatCompletionAssistantMessage("好的，我会尽力办好群助理的角色，请问我有什么需要我帮助的吗？"), s.Prompt())
	}
	if config.C.ContextStatus {
		// 只有在用户开启上下文的时候，追加上下文需要加锁,得到回复追加上下文后才进行锁的释放
		um.Lock()
		defer um.Unlock()
		ctxMsg := um.buildMessage(sender, currentMsg)
		msgs = append(msgs, ctxMsg...)

	} else {
		msgs = append(msgs, currentMsg)
	}
	//fmt.Println("==========")
	//fmt.Println(len(msgs))
	//fmt.Println(msgs)
	req := &qianfan.ChatCompletionRequest{
		Messages: msgs,
	}
	// 发送消息
	reply, err := s.client.Do(ctx, req)
	if nil != err {
		// 发送失败嘞
		return err.Error()
	}
	// 发送成功，可以将请求和回复加入上下文
	if config.C.ContextStatus {
		// 如果请求成功才把问题回复都添加进上下文
		if resetMsg, ok := plugins.Manger.DoPlugin(reply.Result); ok {
			reply.Result = resetMsg
			goto RETURN
		}
		um.addContext(currentMsg, s.Prompt())
		um.addContext(
			qianfan.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: reply.Result}, s.prompt)
	}
RETURN:
	return reply.Result
}

func (s *baiduSpeedSession) CreateImage(ctx context.Context, prompt string) string {
	return ""
}

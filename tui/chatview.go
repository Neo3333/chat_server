package tui
// ref: https://github.com/marcusolsson/tui-go/blob/master/example/chat/main.go

import(
	"errors"
	"fmt"
	tui "github.com/marcusolsson/tui-go"
	"net"
	"time"
)
const(
	MESSAGE = "Message: "
	RECEIVER = "TO: "
)

type SubmitMessageHandler func(string,string)

var chatChain *tui.SimpleFocusChain

type ChatView struct {
	tui.Box
	public   		bool
	name 			*tui.StatusBar
	frame    		*tui.Box
	history  		*tui.Box
	onSubmit 		SubmitMessageHandler
}

type post struct {
	username string
	message  string
	time     string
}

var posts = []post{
	{username: "system", message: "Welcome to Skynet!", time: "00:00"},
}

func getClientIp() (string ,error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "",err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(),nil
			}

		}
	}
	return "", errors.New("Can not find the client ip address!")
}

func NewChatView() *ChatView {
	view := &ChatView{public: true,name: tui.NewStatusBar("")}
	chatChain = &tui.SimpleFocusChain{}

	status := tui.NewStatusBar("MMode: Public.")
	change := tui.NewButton("[Change Mode]")
	change.OnActivated(func(b *tui.Button) {
		if view.public{
			status.SetText("MMode: Private.")
			view.public = false
		}else {
			status.SetText("MMode: Public.")
			view.public = true
		}
	})

	ip, err := getClientIp()
	if err != nil{
		ip = "127.0.0.1"
	}

	sidebar := tui.NewVBox(
		tui.NewLabel("SKYNET COMMUNICATION"),
		tui.NewLabel(""),
		tui.NewLabel(fmt.Sprintf("IP: %v",ip)),
		view.name,
		tui.NewSpacer(),
		change,
	)
	sidebar.SetBorder(true)

	view.history = tui.NewVBox()
	for _, m := range posts {
		view.history.Append(tui.NewHBox(
			tui.NewLabel(time.Now().Format("2006-01-02 15:04:05")),
			tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", m.username))),
			tui.NewLabel(m.message),
			tui.NewSpacer(),
		))
	}
	historyScroll := tui.NewScrollArea(view.history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	input.SetText(MESSAGE)

	receive := tui.NewEntry()
	receive.SetSizePolicy(tui.Minimum,tui.Maximum)
	receive.SetText(RECEIVER)

	input.OnSubmit(func(e *tui.Entry) {
		if view.public{
			if e.Text() != "" {
				if view.onSubmit != nil {
					view.onSubmit(e.Text()[len(MESSAGE):],"")
				}
				e.SetText(MESSAGE)
			}
		}else{
			if e.Text() != "" && receive.Text() != ""{
				if view.onSubmit != nil{
					view.onSubmit(e.Text()[len(MESSAGE):],
						receive.Text()[len(RECEIVER):])
				}
				e.SetText(MESSAGE)
				receive.SetText(RECEIVER)
			}
		}
	})

	inputBox := tui.NewHBox(receive,input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox,inputBox)
	chat.SetSizePolicy(tui.Expanding,tui.Expanding)

	content := tui.NewHBox(sidebar,chat,)
	root := tui.NewVBox(content,status,)

	view.frame = root
	view.frame.SetBorder(true)
	view.Append(view.frame)

	chatChain.Set(input,change,receive)
	return view
}

func (c *ChatView) OnSubmit(handler SubmitMessageHandler) {
	c.onSubmit = handler
}

func (c *ChatView) AddMessage(user string, msg string, time string) {
	c.history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%v <%v> %v", time, user, msg)),
		),
	)
}

func (c *ChatView) SetName(name string)  {
	c.name.SetText(fmt.Sprintf("Username: %s",name))
}


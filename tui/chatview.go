package tui
// ref: https://github.com/marcusolsson/tui-go/blob/master/example/chat/main.go

import(
	"fmt"
	tui "github.com/marcusolsson/tui-go"
	"time"
)

type SubmitMessageHandler func(string)

type ChatView struct {
	tui.Box
	frame    *tui.Box
	history  *tui.Box
	onSubmit SubmitMessageHandler
}

type post struct {
	username string
	message  string
	time     string
}

var posts = []post{
	{username: "system", message: "Welcome to Skynet!", time: "00:00"},
}

func NewChatView() *ChatView {
	view := &ChatView{}
	sidebar := tui.NewVBox(
		tui.NewLabel("CHANNELS"),
		tui.NewLabel("general"),
		tui.NewLabel("random"),
		tui.NewLabel(""),
		tui.NewLabel("DIRECT MESSAGES"),
		tui.NewLabel("slackbot"),
		tui.NewSpacer(),
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

	input.OnSubmit(func(e *tui.Entry) {
		if e.Text() != "" {
			if view.onSubmit != nil {
				view.onSubmit(e.Text())
			}

			e.SetText("")
		}
	})

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox,inputBox)
	chat.SetSizePolicy(tui.Expanding,tui.Expanding)

	root := tui.NewHBox(chat)

	view.frame = root
	view.frame.SetBorder(true)
	view.Append(view.frame)

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

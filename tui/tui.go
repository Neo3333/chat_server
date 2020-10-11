package tui

import (
	"../client"
	"github.com/marcusolsson/tui-go"
)

func StartUi(c client.ChatClient) {

	loginView := NewLoginView()
	chatView := NewChatView()

	ui, err := tui.New(loginView)
	if err != nil {
		panic(err)
	}

	quit := func() { ui.Quit() }

	ui.SetKeybinding("Esc", quit)
	ui.SetKeybinding("Ctrl+c", quit)
	ui.SetFocusChain(loginChain)

	loginView.OnLogin(func(username string) {
		_ = c.SetName(username)
		ui.SetWidget(chatView)
		ui.SetFocusChain(chatChain)
		loginChain = nil
	})

	chatView.OnSubmit(func(msg string) {
		_ = c.SendMessage(msg)
	})
	chatView.OnPrivate(func(msg string, rec string) {
		_ = c.SendMessagePrivate(msg,rec)
	})

	go func() {
		for{
			select {
			case msg := <-c.Incoming():
				ui.Update(func() {
					chatView.AddMessage(msg.Name, msg.Message, msg.Time)
				})
			case err := <-c.Errors():
				ui.Update(func() {
					chatView.AddMessage("system", err.Message, err.Time)
				})
			}
		}
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}


package tui

import (
	"../client"
	"github.com/marcusolsson/tui-go"
)

/**
version 1.0
*/

var ch = make(chan struct{})
var ptr *string

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

		ptr = &username
		ch <- struct{}{}
	})

	chatView.OnSubmit(func(msg string, rec string) {
		_ = c.SendMessage(msg,rec)
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
			case <-c.Done():
				break
			}
		}
	}()

	go func() {
		<-ch
		ui.Update(func() {
			chatView.SetName(*ptr)
		})
		close(ch)
		ptr = nil
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

